#!/bin/bash

set -e

function usage () {
  echo "Usage: "
  echo "  runtest.sh [-C channelName] [-c chaincodename]"
  echo "  runtest.sh -h|--help (print this message)"
  echo "     -C channel name"
  echo "     -c chaincode name"
  echo
  exit 1
}

# Parse commandline args
while getopts "h?C:c:" opt; do
  case "$opt" in
    h|\?)
      usage
      exit 1
    ;;
    C)  TOTAL_CH="$OPTARG"
    ;;
    c)  TOTAL_CC="$OPTARG"
    ;;
  esac
done

: ${TOTAL_CH:=2}
: ${TOTAL_CC:=4}
: ${CHANNEL_NAME:="ch"}

for (( i=1;i<=$TOTAL_CC;i=$i+1 ))
do
	peer chaincode install -n mycc$i -v 0 -p  github.com/hyperledger/fabric/examples/ccchecker/chaincodes/newkeyperinvoke/
done

for (( i=1;i<=$TOTAL_CH;i=$i+1 ))
do
	printf "\n\n *********** Creating channel configuration '$CHANNEL_NAME$i.tx'... *********** \n\n"
	configtxgen -channelID $CHANNEL_NAME$i -outputCreateChannelTx $CHANNEL_NAME$i.tx -profile SampleSingleMSPChannel 

	printf "\n\n *********** Creating channel '$CHANNEL_NAME$i' ... *********** \n\n"
	peer channel create -o 127.0.0.1:7050 -c $CHANNEL_NAME$i -f $CHANNEL_NAME$i.tx -t 10

	printf "\n\n *********** Joining the peer on channel '$CHANNEL_NAME$i' *************\n\n"
	peer channel join -b $CHANNEL_NAME$i.block
done

for (( ch=1;ch<=$TOTAL_CH;ch=$ch+1 ))
do
	for (( cc=1;cc<=$TOTAL_CC;cc=$cc+1 ))
	do
		printf "\n\n *********** Instantiate chaincode mycc$cc on channel '$CHANNEL_NAME$ch' *************\n\n"
		peer chaincode instantiate -n mycc$cc -v 0 -c '{"Args":[""]}' -o 127.0.0.1:7050 -C $CHANNEL_NAME$ch
				
	done
done

sleep 20

(./ccchecker -e env.json -c ccchecker1.json > log1.txt ) &
(./ccchecker -e env.json -c ccchecker2.json > log2.txt ) &

printf "\n\n ******* Execution started monitor log1.txt & log2.txt files ********* \n\n"


