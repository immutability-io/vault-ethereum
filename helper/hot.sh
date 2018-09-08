#!/bin/bash


function print_help {
    echo "Usage: bash hot.sh OPTIONS"
    echo -e "\nOPTIONS:"
    echo -e "  [keybase]\tName of Keybase user used to encrypt Vault keys"
    echo -e "  [path]\tPath to cold storage"
}

if [ -z "$2" ]; then
    print_help
    exit 0
elif [ "$1" == "--help" ]; then
    print_help
    exit 0
else
  KEYBASE_USER=$1
  COLD_STORAGE=$2
fi

cp -r $COLD_STORAGE/etc $HOME/etc
unset VAULT_TOKEN
nohup /usr/local/bin/vault server -config $HOME/etc/vault.d/vault.hcl &> /dev/null &

cp $COLD_STORAGE/"$KEYBASE_USER"_* .
vault operator unseal $(keybase decrypt -i $KEYBASE_USER"_UNSEAL_1.txt")
vault operator unseal $(keybase decrypt -i $KEYBASE_USER"_UNSEAL_2.txt")
vault operator unseal $(keybase decrypt -i $KEYBASE_USER"_UNSEAL_3.txt")
