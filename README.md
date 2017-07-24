# ChaincodeChecker (ccc) using native binaries ( peer & orderer)

This repository is to maintain the latest chaincodechecker code (This is temporary untill the changes been merged into Fabric/examples/ccchecker)

### How to run ccc:

* Generate/Use the peer/orderer binaries from fabric repo. 
* Clone the repo 
  git clone https://github.com/asararatnakar/ccc.git

- execute the following commands from different terminals :
* cd ccc

### Terminal-1

```
mkdir -p hyperledger/production/orderer

export PATH=$PATH:$PWD/bin

export FABRIC_CFG_PATH=$PWD/sampleconfig

ORDERER_FILELEDGER_LOCATION=./hyperledger/production/orderer ORDERER_GENERAL_GENESISPROFILE=SampleSingleMSPSolo orderer
```

### Terminal-2

```
export PATH=$PATH:$PWD/bin

export FABRIC_CFG_PATH=$PWD/sampleconfig

CORE_PEER_FILESYSTEMPATH=./hyperledger/production peer node start

```

### Terminal-3

```
export PATH=$PATH:$PWD/bin

export FABRIC_CFG_PATH=$PWD/sampleconfig

./runtest.sh
```
 execution of the shell script `runtest.sh` results to creation of channel artifacts , channel creation,  join channel etc.,. 
  also installs and instantiates the chaincode.

 This defaults to 2 channels and 4 chaincodes. the defaults can be overridden by providing the followingf options 
  ` -C `  - # of Channels
  ` -c `  - # of chaincodes

```
	ex: ./runtest.sh -C 3 -c 10
```

** 3 Channels and 10 chaincode s**

* Script continues to start two channels, each channel with # of chaincodes specified.
  Make sure you change the ccchecker<N>.json files as per your inputs to the script
