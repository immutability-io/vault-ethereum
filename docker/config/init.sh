#!/bin/bash

OPERATOR_JSON="/home/vault/config/operator.json"
OPERATOR_SECRETS=$(cat $OPERATOR_JSON)


function banner() {
  echo "+----------------------------------------------------------------------------------+"
  printf "| %-80s |\n" "`date`"
  echo "|                                                                                  |"
  printf "| %-80s |\n" "$@"
  echo "+----------------------------------------------------------------------------------+"
}

function authenticate() {
    banner "Authenticating to $VAULT_ADDR as root"
    ROOT=$(echo $OPERATOR_SECRETS | jq -r .root_token)
    export VAULT_TOKEN=$ROOT
}

function unauthenticate() {
    banner "Unsetting VAULT_TOKEN"
    unset VAULT_TOKEN
}

function unseal() {
    banner "Unsealing $VAULT_ADDR..."
    UNSEAL=$(echo $OPERATOR_SECRETS | jq -r '.unseal_keys_hex[0]')
    vault operator unseal $UNSEAL
}

function configure() {
    banner "Installing vault-ethereum plugin at $VAULT_ADDR..."
	SHA256SUMS=`cat /home/vault/plugins/SHA256SUMS | awk '{print $1}'`
	vault write sys/plugins/catalog/secret/vault-ethereum \
		  sha_256="$SHA256SUMS" \
		  command="vault-ethereum --ca-cert=$CA_CERT --client-cert=$TLS_CERT --client-key=$TLS_KEY"

	if [[ $? -eq 2 ]] ; then
	  echo "vault-ethereum couldn't be written to the catalog!"
	  exit 2
	fi

	vault secrets enable -path=vault-ethereum -plugin-name=vault-ethereum plugin
	if [[ $? -eq 2 ]] ; then
	  echo "vault-ethereum couldn't be enabled!"
	  exit 2
	fi
    vault audit enable file file_path=stdout
}

function status() {
    vault status
}

function init() {
    OPERATOR_SECRETS=$(vault operator init -key-shares=1 -key-threshold=1 -format=json | jq .)
    echo $OPERATOR_SECRETS > $OPERATOR_JSON
}
sleep 20
if [ -f "$OPERATOR_JSON" ]; then
    unseal
    status
else
    init
    unseal
    authenticate
    configure
    unauthenticate
    status
fi
