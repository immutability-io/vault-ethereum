#!/bin/bash


function print_help {
    echo "Usage: bash hot.sh OPTIONS"
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

vault operator unseal $(keybase decrypt -i $KEYBASE_USER"_UNSEAL_1.txt")
vault operator unseal $(keybase decrypt -i $KEYBASE_USER"_UNSEAL_2.txt")
vault operator unseal $(keybase decrypt -i $KEYBASE_USER"_UNSEAL_3.txt")
