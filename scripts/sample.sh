#!/bin/bash

# Set localnet configuration
# Reference localnet script to see which tokens are given to the user accounts in genesis state
BINARY=farmingd
CHAIN_ID=localnet
CHAIN_DIR=./data
USER_1_ADDRESS=cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu
USER_2_ADDRESS=cosmos185fflsvwrz0cx46w6qada7mdy92m6kx4gqx0ny

# Ensure jq is installed
if [[ ! -x "$(which jq)" ]]; then
  echo "jq (a tool for parsing json in the command line) is required..."
  echo "https://stedolan.github.io/jq/download/"
  exit 1
fi

# Ensure farmingd is installed
if ! [ -x "$(which $BINARY)" ]; then
  echo "Error: $BINARY is not installed. Try building $BINARY by 'make install'" >&2
  exit 1
fi

# Ensure localnet is running
if [[ "$(pgrep $BINARY)" == "" ]];then
    echo "Error: localnet is not running. Try running localnet by 'make localnet" 
    exit 1
fi
