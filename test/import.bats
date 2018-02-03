#!/usr/bin/env bats

@test "import test account" {
  passphrase=$(cat passphrase.txt)
  import_path=$(cat path.txt)
  import_address="$(vault write -format=json ethereum/import/test2 chain_id=1977 passphrase=$passphrase path=$import_path | jq .data.address | tr -d '"')"
  export_address="$(vault read -format=json ethereum/import/test | jq .data.address | tr -d '"')"
    [ "$import_address" != "$export_address" ]
}
