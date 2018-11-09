#!/bin/bash
#
# SPDX-License-Identifier: Apache-2.0
# This code is based on code written by the Hyperledger Fabric community. 
# Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/fabcar/startFabric.sh
#
# Exit on first error

set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

CURRENT_DIR=$PWD
starttime=$(date +%s)

if [ ! -d ~/.hfc-key-store/ ]; then
	mkdir ~/.hfc-key-store/
fi

# launch network; create channel and join peer to channel
cd ../basic-network
./start.sh

# Now launch the CLI container in order to install, instantiate chaincode
# and prime the ledger with our 10 tuna catches
docker-compose -f ./docker-compose.yml up -d cli

docker exec cli peer chaincode install -n con-app -v 1.5 -p github.com/mycontract
docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n con-app -v 1.5 -c '{"Args":[""]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
sleep 10
docker exec cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n con-app -c '{"function":"initLedger","Args":[""]}'

sleep 10
printf "\nStart node container\n"


cd "$CURRENT_DIR"
docker-compose -f ./docker-compose-node.yml up
printf "\nTotal execution time : $(($(date +%s) - starttime)) secs ...\n\n"
#printf "\nStart with the enrollAdmin.js, then registerUser.js, then server.js\n\n"
