#!/usr/bin/env bats

@test "test create accounts" {
  testdata=$(cat ./_NAME_/names.json)
  for data in $(echo "${testdata}" | jq -r '.[]'); do
    account="$(vault write -format=json -f ethereum/accounts/$data | jq .data)"
    address=$(echo $account | jq .address)
      [ "${#address}" -eq 44 ]

  done
}

@test "test list accounts" {
  testdata=$(cat ./_NAME_/names.json)
  testdata_sha=$(echo $testdata | shasum -a 256)
  testdata2="$(vault list -format=json ethereum/accounts)"
  testdata2_sha=$(echo $testdata2 | shasum -a 256)
    [ "$testdata_sha" = "$testdata2_sha" ]
}

@test "test delete accounts" {
  testdata=$(cat ./_NAME_/names.json)
  for data in $(echo "${testdata}" | jq -r '.[]'); do
    account="$(vault delete ethereum/accounts/$data)"
  done
}
