#!/bin/bash

function print_help {
    echo "Usage: bash install.sh OPTIONS"
    echo -e "\nOPTIONS:"
    echo -e "  [keybase]\tName of Keybase user to encrypt Vault keys with"
}

function initialize {
  export VAULT_ADDR=https://localhost:8200
  export VAULT_CACERT=$HOME/etc/vault.d/root.crt
  export VAULT_INIT=$(vault operator init -format=json)
  if [[ $? -eq 2 ]] ; then
    echo "Vault initialization failed!"
    exit 2
  fi
  ROOT_TOKEN=$(echo $VAULT_INIT | jq -r .root_token)
  keybase encrypt $KEYBASE -m $ROOT_TOKEN -o ./"$KEYBASE"_VAULT_ROOT_TOKEN.txt
  if [[ $? -eq 2 ]] ; then
    echo "Keybase encryption failed!"
    exit 2
  fi
  for (( COUNTER=0; COUNTER<5; COUNTER++ ))
  do
    key=$(echo $VAULT_INIT | jq -r '.unseal_keys_hex['"$COUNTER"']')
    vault operator unseal $key
    keybase encrypt $KEYBASE -m $key -o ./"$KEYBASE"_UNSEAL_"$COUNTER".txt
  done
  unset VAULT_INIT
  unset ROOT_TOKEN
}

if [ -z "$1" ]; then
    print_help
    exit 0
else
    KEYBASE=$1
fi

unset VAULT_TOKEN
unset VAULT_ADDR
unset VAULT_CACERT

nohup /usr/local/bin/vault server -config $HOME/etc/vault.d/vault.hcl &> /dev/null &
sleep 10

initialize

