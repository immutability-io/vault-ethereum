#!/bin/bash

function install_plugin {
  echo "ADDING TO CATALOG: sys/plugins/catalog/ethereum-plugin"
  vault write sys/plugins/catalog/ethereum-plugin \
        sha_256="$(cat SHA256SUM)" \
        command="vault-ethereum --ca-cert=$HOME/etc/vault.d/root.crt --client-cert=$HOME/etc/vault.d/vault.crt --client-key=$HOME/etc/vault.d/vault.key"

  if [[ $? -eq 2 ]] ; then
    echo "Vault Catalog update failed!"
    exit 2
  fi

  echo "MOUNTING: ethereum/mainnet"
  vault secrets enable -path=ethereum/mainnet -description="Immutability's Ethereum Wallet - Mainnet" -plugin-name=ethereum-plugin plugin
  if [[ $? -eq 2 ]] ; then
    echo "Failed to mount Ethereum plugin for mainnet!"
    exit 2
  fi
  echo "MOUNTING: ethereum/rinkeby"
  vault secrets enable -path=ethereum/rinkeby -description="Immutability's Ethereum Wallet - Rinkeby" -plugin-name=ethereum-plugin plugin
  if [[ $? -eq 2 ]] ; then
    echo "Failed to mount Ethereum plugin for rinkeby!"
    exit 2
  fi
  echo "CONFIGURE: ethereum/mainnet"
  vault write ethereum/mainnet/config rpc_url="https://mainnet.infura.io" chain_id="1"
  echo "CONFIGURE: ethereum/rinkeby"
  vault write -f ethereum/rinkeby/config
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

export VAULT_TOKEN=$(keybase decrypt -i $KEYBASE_USER"_VAULT_ROOT_TOKEN.txt")

install_plugin

unset VAULT_TOKEN