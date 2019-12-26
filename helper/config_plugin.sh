#!/bin/bash

function install_plugin {
  echo "ADDING TO CATALOG: sys/plugins/catalog/secret/ethereum-plugin"
  vault write sys/plugins/catalog/secret/ethereum-plugin \
        sha_256="$(cat SHA256SUM)" \
        command="vault-ethereum --ca-cert=$HOME/etc/vault.d/root.crt --client-cert=$HOME/etc/vault.d/vault.crt --client-key=$HOME/etc/vault.d/vault.key"

  if [[ $? -eq 2 ]] ; then
    echo "Vault Catalog update failed!"
    exit 2
  fi

  echo "MOUNTING: ethereum/prod"
  vault secrets enable -path=ethereum/prod -description="Immutability's Ethereum Wallet - PROD" -plugin-name=ethereum-plugin plugin
  if [[ $? -eq 2 ]] ; then
    echo "Failed to mount Ethereum plugin for prod!"
    exit 2
  fi
  echo "MOUNTING: ethereum/dev"
  vault secrets enable -path=ethereum/dev -description="Immutability's Ethereum Wallet - DEV" -plugin-name=ethereum-plugin plugin
  if [[ $? -eq 2 ]] ; then
    echo "Failed to mount Ethereum plugin for dev!"
    exit 2
  fi
  echo "CONFIGURE: ethereum/prod"
  vault write ethereum/prod/config rpc_url="https://mainnet.infura.io" chain_id="1"
  echo "CONFIGURE: ethereum/dev"
  vault write -f ethereum/dev/config
}

function print_help {
    echo "Usage: bash config_plugin.sh OPTIONS"
    echo -e "\nOPTIONS:"
    echo -e "  [keybase]\tName of Keybase user used to encrypt Vault keys"
}

if [ -z "$1" ]; then
    print_help
    exit 0
elif [ "$1" == "--help" ]; then
    print_help
    exit 0
else
  KEYBASE_USER=$1
fi

source ./.as-root $KEYBASE_USER

install_plugin

unset VAULT_TOKEN