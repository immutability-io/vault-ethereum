#!/usr/bin/env bats

@test "test write rinkeby config" {
  config="$(vault write -format=json ethereum/dev/config rpc_url="https://rinkeby.infura.io" chain_id="4" | jq .data)"
  api_key="$(echo $config | jq -r .api_key)"
  bound_cidr_list="$(echo $config | jq -r .bound_cidr_list)"
  chain_id="$(echo $config | jq -r .chain_id)"
  rpc_url="$(echo $config | jq -r .rpc_url)"
    [ "$api_key" = "" ]
    [ "$chain_id" = "4" ]
    [ "$rpc_url" = "https://rinkeby.infura.io" ]
}

@test "test write mainnet config" {
  config="$(vault write -format=json ethereum/prod/config rpc_url="https://mainnet.infura.io" chain_id="1" api_key="$COINMARKETCAP_API_KEY"  | jq .data)"
  api_key="$(echo $config | jq -r .api_key)"
  bound_cidr_list="$(echo $config | jq -r .bound_cidr_list)"
  chain_id="$(echo $config | jq -r .chain_id)"
  rpc_url="$(echo $config | jq -r .rpc_url)"
    [ "$api_key" = "$COINMARKETCAP_API_KEY" ]
    [ "$chain_id" = "1" ]
    [ "$rpc_url" = "https://mainnet.infura.io" ]
}

@test "test read config" {
  config="$(vault write -format=json ethereum/prod/config rpc_url="https://mainnet.infura.io" chain_id="1" api_key="$COINMARKETCAP_API_KEY"  | jq .data)"
  api_key="$(echo $config | jq -r .api_key)"
  bound_cidr_list="$(echo $config | jq -r .bound_cidr_list)"
  chain_id="$(echo $config | jq -r .chain_id)"
  rpc_url="$(echo $config | jq -r .rpc_url)"
  read_config="$(vault read -format=json ethereum/prod/config | jq .data)"
  read_api_key="$(echo $read_config | jq -r .api_key)"
  read_bound_cidr_list="$(echo $read_config | jq -r .bound_cidr_list)"
  read_chain_id="$(echo $read_config | jq -r .chain_id)"
  read_rpc_url="$(echo $read_config | jq -r .rpc_url)"
    [ "$api_key" = "$read_api_key" ]
    [ "$chain_id" = "$read_chain_id" ]
    [ "$rpc_url" = "$read_rpc_url" ]
}
