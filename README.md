# ChaincodeChecker (ccc) using native binaries ( peer & orderer)

This repository is to maintain the latest chaincodechecker code (This is temporary untill the changes been merged into Fabric/examples/ccchecker)

### How to run ccc:

* Generate/Use the peer/orderer binaries from fabric repo. 
* Clone the repo 
  git clone https://github.com/asararatnakar/ccc.git

- execute the commands:
* cd ccc
* execute the shell script `runtest.sh` to create channel artifacts , create channel,  join channel. also install and instantiates the chaincode.
  `./runtest.sh`

 This defaults to 2 channels and 4 chaincodes. the defaults can be overridden by providing the followingf options 
  ` -C `  - # of Channels
  ` -c `  - # of chaincodes

	```
	 ex: ./runtest.sh -C 3 -c 10
	```

	here we are asking for 3 Channels and 10 chaincodes

* Script continues to start two channels, each channel with # of chaincodes specified.
  Make sure you change the ccchecker<N>.json files as per your inputs to the script


