#!/usr/bin/env bats

@test "test write config" {
  config="$(vault write -format=json -f ethereum/config | jq .data)"
  api_key="$(echo $config | jq -r .api_key)"
  bound_cidr_list="$(echo $config | jq -r .bound_cidr_list)"
  chain_id="$(echo $config | jq -r .chain_id)"
  rpc_url="$(echo $config | jq -r .rpc_url)"
    [ "$api_key" = "" ]
    [ "$chain_id" = "4" ]
    [ "$rpc_url" = "https://rinkeby.infura.io" ]
}

@test "test read config" {
  config="$(vault write -format=json -f ethereum/config | jq .data)"
  api_key="$(echo $config | jq -r .api_key)"
  bound_cidr_list="$(echo $config | jq -r .bound_cidr_list)"
  chain_id="$(echo $config | jq -r .chain_id)"
  rpc_url="$(echo $config | jq -r .rpc_url)"
  read_config="$(vault read -format=json ethereum/config | jq .data)"
  read_api_key="$(echo $read_config | jq -r .api_key)"
  read_bound_cidr_list="$(echo $read_config | jq -r .bound_cidr_list)"
  read_chain_id="$(echo $read_config | jq -r .chain_id)"
  read_rpc_url="$(echo $read_config | jq -r .rpc_url)"
    [ "$api_key" = "$read_api_key" ]
    [ "$chain_id" = "$read_chain_id" ]
    [ "$rpc_url" = "$read_rpc_url" ]
}
