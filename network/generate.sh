#!/bin/bash

set -x

export PATH=$PATH:/home/bstudent/fabric-samples/bin;
export FABRIC_CFG_PATH=${PWD}

if [ ! -d config ]; then
    config
fi 

rm -rf ./config/*
rm -rf ./crypto-config

#1 crypto generation
cryptogen generate --config=./crypto-config.yaml 

#2 genesis block generation
configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock ./config/genesis.block

#3 channel transaction generation
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./config/mychannel.tx -channelID mychannel

#4 anchor peer transaction generation
