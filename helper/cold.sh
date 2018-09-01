#!/bin/bash


function print_help {
    echo "Usage: bash cold.sh OPTIONS"
    echo -e "\nOPTIONS:"
    echo -e "  [keybase]\tName of Keybase user used to encrypt Vault keys"
    echo -e "  [path]\tPath to mounted Flash drive or other media"
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

kill -2 $(ps aux | grep '/usr/local/bin/vault server' | awk '{print $2}')
mv -f $HOME/etc $COLD_STORAGE/etc
mv "$KEYBASE_USER"_* $COLD_STORAGE
