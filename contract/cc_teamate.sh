#!/bin/bash

CCNAME=teamate
CC_SRC_PATH=github.com/teamate
CC_VERSION=1.0.0
CC_SEQUENCE=1
CC_END_POLICY="NA"

if [ "$CC_END_POLICY" = "NA" ]; then
  CC_END_POLICY=""
else
  CC_END_POLICY="--signature-policy $CC_END_POLICY"
fi

PEER_CONN_PARMS="--peerAddresses localhost:7051"


FABRIC_CFG_PATH=/home/bstudent/fabric-samples/config/

setGlobals() {
    USING_ORG=$1
    if [ $USING_ORG -eq 1 ]; then
        export CORE_PEER_LOCALMSPID="Org1MSP"
        export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
        export CORE_PEER_ADDRESS=localhost:7051
    elif [ $USING_ORG -eq 2 ]; then
        export CORE_PEER_LOCALMSPID="Org2MSP"
        export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
        export CORE_PEER_ADDRESS=localhost:9051
    else
        errorln "ORG Unknown"
    fi
}

peer lifecycle chaincode package ${CC_NAME}.tar.gz --path ${CC_SRC_PATH} --label ${CC_NAME}_${CC_VERSION} 

setGlobals 1
peer lifecycle chaincode install ${CC_NAME}.tar.gz 
peer lifecycle chaincode queryinstalled 

setGlobals 2
peer lifecycle chaincode install ${CC_NAME}.tar.gz 
peer lifecycle chaincode queryinstalled 

setGlobals 1
peer lifecycle chaincode approveformyorg -o localhost:7050 --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} ${CC_END_POLICY} 

peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${CC_END_POLICY} --output json 

setGlobals 2
peer lifecycle chaincode approveformyorg -o localhost:7050 --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} ${CC_END_POLICY} 

peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${CC_END_POLICY} --output json 

setGlobals 1
peer lifecycle chaincode commit -o localhost:7050 --channelID $CHANNEL_NAME --name ${CC_NAME} $PEER_CONN_PARMS --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${CC_END_POLICY} 
peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name ${CC_NAME} 


fcn_call='{"function":"addUser","Args":["user1"]}'
peer chaincode invoke -o localhost:7050 -C $CHANNEL_NAME -n ${CC_NAME} --isInit -c ${fcn_call} 

peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"Args":["readRating"]}' 


