#!/bin/bash

COMMAND=$1
PWD="$(pwd)"
GANACHE_DATA="$PWD/ganache_data"
OPERATOR_JSON="$PWD/config/operator.json"
OPERATOR_DATA="$PWD/config/data"
OPERATOR_SECRETS=$(cat $OPERATOR_JSON)
export VAULT_CACERT="$PWD/config/root.crt"
export VAULT_ADDR=https://localhost:9200


function clear() {
    echo "Clearing previous state..."
    echo "rm -fr $GANACHE_DATA"
    rm -fr $GANACHE_DATA
    echo "rm -fr $OPERATOR_DATA"
    rm -fr $OPERATOR_DATA
    echo "rm $OPERATOR_JSON"
    rm $OPERATOR_JSON
    rm $PWD/config/*.key
    rm $PWD/config/*.crt
    rm $PWD/config/*.cnf
    rm $PWD/config/*.srl
    rm $PWD/config/*.csr
    rm $PWD/README.md
    
}

function authenticate() {
    echo "Authenticating to $VAULT_ADDR as root"
    ROOT=$(echo $OPERATOR_SECRETS | jq -r .root_token)
    export VAULT_TOKEN=$ROOT
}

function unseal() {
    echo "Unsealing $VAULT_ADDR..."
    UNSEAL=$(echo $OPERATOR_SECRETS | jq -r '.unseal_keys_hex[0]')
    vault operator unseal $UNSEAL
}

function status() {
    vault status
    vault secrets list
}

function init() {
    OPERATOR_SECRETS=$(vault operator init -key-shares=1 -key-threshold=1 -format=json | jq .)
    echo $OPERATOR_SECRETS > $OPERATOR_JSON
}

if [ "$#" -ne 1 ]; then
    echo "Illegal number of parameters"
    exit 1
fi

if [ $COMMAND = "auth" ]; then
    authenticate
elif [ $COMMAND = "unseal" ]; then
    unseal
elif [ $COMMAND = "clear" ]; then
    clear
fi

