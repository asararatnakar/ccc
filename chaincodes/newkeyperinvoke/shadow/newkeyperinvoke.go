/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shadow

import (
	"bytes"
	"fmt"
	"sync"
)

// NewKeyPerInvoke is the shadow implementation for  NewKeyPerInvoke in the parent package
// The shadow provides invoke arguments that are guaranteed to be result in unique ledger
// entries as long as the parameters to GetInvokeArgs are unique
type NewKeyPerInvoke struct {
	sync.Mutex
	state map[string][]byte

	//iterator to be used in validation phase for the "Next" call
	iterator *stateIterator
}

//---------- implements ShadowCCIntf functions -------

//Clone yourself
func (t *NewKeyPerInvoke) Clone() interface{} {
	return &NewKeyPerInvoke{}
}

//InitShadowCC initializes CC
func (t *NewKeyPerInvoke) InitShadowCC(ccname string, initArgs []string) {
	t.state = make(map[string][]byte)
}

//invokeSuccessful sets the state and increments succefull invokes counter
func (t *NewKeyPerInvoke) invokeSuccessful(key []byte, val []byte) {
	t.Lock()
	defer t.Unlock()
	t.state[string(key)] = val
}

//getState gets the state
func (t *NewKeyPerInvoke) getState(key []byte) ([]byte, bool) {
	t.Lock()
	defer t.Unlock()
	v, ok := t.state[string(key)]
	return v, ok
}

//OverrideNumInvokes returns the number of invokes shadow wants
//accept users request, no override
func (t *NewKeyPerInvoke) OverrideNumInvokes(numInvokesPlanned int) int {
	return numInvokesPlanned
}

//GetNumQueries returns the number of queries shadow wants ccchecked to do.
//For our purpose, just do as many queries as there were invokes for.
func (t *NewKeyPerInvoke) GetNumQueries(numInvokesCompletedSuccessfully int) int {
	return numInvokesCompletedSuccessfully
}

//GetInvokeArgs get args for invoke based on chaincode ID and iteration num
func (t *NewKeyPerInvoke) GetInvokeArgs(ccnum int, iter int) [][]byte {
	args := make([][]byte, 3)
	args[0] = []byte("put")
	args[1] = []byte(fmt.Sprintf("%d_%d", ccnum, iter))
	args[2] = []byte(fmt.Sprintf("%d", ccnum))

	return args
}

//PostInvoke store the the key/val for later verification
func (t *NewKeyPerInvoke) PostInvoke(args [][]byte, resp []byte) error {
	if len(args) < 3 {
		return fmt.Errorf("invalid number of args posted %d", len(args))
	}

	if string(args[0]) != "put" {
		return fmt.Errorf("invalid args posted %s", args[0])
	}

	//the actual CC should have returned OK for success
	if string(resp) != "OK" {
		return fmt.Errorf("invalid response %s", string(resp))
	}

	t.invokeSuccessful(args[1], args[2])

	return nil
}

// -------- validation phase ----------

type stateIterator struct {
	nkpi  *NewKeyPerInvoke
	argsC chan [][]byte
}

//getQueryArgs returns the query for ccchecker to test against
func (i *stateIterator) getQueryArgs(key string) [][]byte {
	args := make([][]byte, 2)
	args[0] = []byte("get")
	args[1] = []byte(key)
	return args
}

func (i *stateIterator) init() {
	//spin off a func to serve on the blocking channel
	//completely reentrant
	go func() {
		defer func() {
			recover()
		}()
		for key := range i.nkpi.state {
			i.argsC <- i.getQueryArgs(key)
		}

		//upto caller to stop calling if nil and close the
		//iterator. This way we completely avoid need for
		//locks and checking. Just use the channel
		for {
			i.argsC <- nil
		}
	}()
}

//can be called concurrently
func (i *stateIterator) next() [][]byte {
	//will be nil if iterator is done
	args := <-i.argsC

	return args
}

func (i *stateIterator) close() {
	close(i.argsC)
}

//InitValidation sets up the shadow chaincode for iteration
//Should be called once
func (t *NewKeyPerInvoke) InitValidation() error {
	t.iterator = &stateIterator{t, make(chan [][]byte)}
	t.iterator.init()
	return nil
}

//NextQueryArgs returns the next args  to call query with
//Could be called any number of times, concurrently
//Could be nil if the iterator is done
func (t *NewKeyPerInvoke) NextQueryArgs() [][]byte {
	if t.iterator == nil {
		return nil
	}
	return t.iterator.next()
}

//ValidationDone cleans up (in this case closes the iterator).
//Should be called once
func (t *NewKeyPerInvoke) ValidationDone() error {
	t.iterator.close()
	return nil
}

//Validate the key/val with mem storage
func (t *NewKeyPerInvoke) Validate(args [][]byte, value []byte) error {
	if len(args) < 2 {
		return fmt.Errorf("invalid number of args for validate %d", len(args))
	}

	if string(args[0]) != "get" {
		return fmt.Errorf("invalid validate function %s", args[0])
	}

	if v, ok := t.getState(args[1]); !ok {
		return fmt.Errorf("key not found %s", args[1])
	} else if !bytes.Equal(v, value) {
		return fmt.Errorf("expected(%s) but found (%s)", string(v), string(value))
	}

	return nil
}
