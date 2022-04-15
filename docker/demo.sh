#!/bin/bash

MNEMONIC="volcano story trust file before member board recycle always draw fiction when"
CHAIN_ID=5777
PORT=8545
RPC_URL="http://ganache:$PORT"
PLUGIN="vault-ethereum"
RAW_TX="f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772"
MESSAGE="HOLA VAULT"
ERC20_CONTRACTS_PATH="$(pwd)/erc20/build/"
CONTRACT_SAFE_MATH="SafeMath"
CONTRACT_OWNED="Owned"
CONTRACT_FIXED_SUPPLY_TOKEN="FixedSupplyToken"
BIN_FILE=".bin"
ABI_FILE=".abi"

function header() {
  printf "## %s\n\n" "$@"
  echo ""
}

function code() {
  echo "\`\`\`"
}

function json() {
  echo "\`\`\`json"
}

function log_vault_command() {
  echo "### Using the Vault CLI"
  echo ""
  code
  printf "%s\n\n" "$@"
  code
  echo ""
}

function log_curl_command() {
  echo "### Using Curl"
  echo ""
  code
  echo "$@"
  code
  echo ""
}

function step() {
  printf "### STEP %s:\n" "$@"
}

function log() {
  echo "### Sample Response"
  printf "\n%s\n" "${date}"
  json
  echo "$@" | jq .
  code
  echo ""
}

function log_json() {
  echo "### Sample Response"
  printf "\n%s\n" "${date}"
  json
  echo "$JSON" | jq .
  code
  echo ""
}

function log_value() {
  echo "### Sample Response"
  printf "\n%s\n" "${date}"
  code
  echo "$@"
  code
  echo ""
}


function vault_command_json() {
  echo "### Using the Vault CLI"
  echo ""
  code
  printf "%s\n\n" "$@"
  code
  JSON=$(eval "$@")
  #echo $JSON | jq .
}

header "HOW TO CONFIGURE THE PLUGIN AFTER IT HAS BEEN ENABLED"
log_curl_command $(vault write -output-curl-string $PLUGIN/config rpc_url="$RPC_URL" chain_id="$CHAIN_ID")
vault_command_json "vault write -format=json $PLUGIN/config  rpc_url='$RPC_URL' chain_id='$CHAIN_ID'"
log_json

header "HOW TO CREATE AN ACCOUNT NAMED BOB USING A MNEMONIC"
vault_command_json "vault write -format=json $PLUGIN/accounts/bob mnemonic='$MNEMONIC'"
log_curl_command $(vault write -output-curl-string $PLUGIN/accounts/bob mnemonic="$MNEMONIC")
log_json

step "1: GET BOB'S ADDRESS"
log_vault_command "vault read -field=address vault-ethereum/accounts/bob"
log_curl_command $(vault read -output-curl-string -field=address vault-ethereum/accounts/bob)
BOB_ADDRESS=$(vault read -field=address vault-ethereum/accounts/bob)
log_value $BOB_ADDRESS

header "HOW TO CREATE AN ACCOUNT NAMED ALICE WITH NO MNEMONIC"
vault_command_json "vault write -f -format=json $PLUGIN/accounts/alice"
log_curl_command $(vault write -f -output-curl-string $PLUGIN/accounts/alice)
log_json

header "HOW TO TRANSFER 0.5 ETH FROM BOB TO ALICE"
step "1: CONVERT 0.5 ETH TO WEI"
log_vault_command "vault write -field=amount_to vault-ethereum/convert unit_from='ETH' unit_to='WEI' amount=0.5"
log_curl_command $(vault write -output-curl-string vault-ethereum/convert unit_from="ETH" unit_to="WEI" amount=0.5)
AMOUNT_TO=$(vault write -field=amount_to vault-ethereum/convert unit_from='ETH' unit_to='WEI' amount=0.5)
log_value $AMOUNT_TO
step "2: GET ALICE'S ADDRESS"
log_vault_command "vault read -field=address vault-ethereum/accounts/alice"
log_curl_command $(vault read -output-curl-string -field=address vault-ethereum/accounts/alice)
ALICE_ADDRESS=$(vault read -field=address vault-ethereum/accounts/alice)
log_value $ALICE_ADDRESS
step "3: SEND BOB'S ETH TO ALICE"
vault_command_json "vault write -format=json vault-ethereum/accounts/bob/transfer to='$ALICE_ADDRESS' amount='$AMOUNT_TO'"
log_curl_command $(vault write -output-curl-string vault-ethereum/accounts/bob/transfer to="$ALICE_ADDRESS" amount="$AMOUNT_TO")
log_json

header "HOW TO CHECK BALANCE OF THE ACCOUNT NAMED BOB"
vault_command_json "vault read -format=json $PLUGIN/accounts/bob/balance"
log_curl_command $(vault read -output-curl-string $PLUGIN/accounts/bob/balance)
log_json

header "HOW TO CHECK BALANCE OF THE ACCOUNT NAMED ALICE"
vault_command_json "vault read -format=json $PLUGIN/accounts/alice/balance"
log_curl_command $(vault read -output-curl-string $PLUGIN/accounts/alice/balance)
log_json

header "HOW TO SIGN A TRANSACTION WITH DATA"
vault_command_json "vault write -format=json $PLUGIN/accounts/alice/sign-tx to='$BOB_ADDRESS' data='$RAW_TX' amount='20000000000000000'"
log_curl_command $(vault write -output-curl-string $PLUGIN/accounts/alice/sign-tx to="$BOB_ADDRESS" data="$RAW_TX" amount='20000000000000000')
log_json

header "HOW TO SIGN MESSAGE. Message signature can be verify on: https://etherscan.io/verifySig/2156" 
vault_command_json "vault write -format=json $PLUGIN/accounts/bob/sign message='$MESSAGE'"
log_curl_command $(vault write -output-curl-string $PLUGIN/accounts/bob/sign message="$MESSAGE")
log_json

header "DEPLOY CONTRACT $CONTRACT_FIXED_SUPPLY_TOKEN"
vault_command_json "vault write -format=json $PLUGIN/accounts/bob/deploy \
abi=@$ERC20_CONTRACTS_PATH$CONTRACT_FIXED_SUPPLY_TOKEN$ABI_FILE \
bin=@$ERC20_CONTRACTS_PATH$CONTRACT_FIXED_SUPPLY_TOKEN$BIN_FILE"

log_json
ERC20_ADDRESS=$(echo $JSON | jq .data.contract)

header "DEPLOY CONTRACT $CONTRACT_OWNED"
vault_command_json "vault write -format=json $PLUGIN/accounts/bob/deploy \
abi=@$ERC20_CONTRACTS_PATH$CONTRACT_OWNED$ABI_FILE \
bin=@$ERC20_CONTRACTS_PATH$CONTRACT_OWNED$BIN_FILE"
log_json

header "DEPLOY CONTRACT $CONTRACT_SAFE_MATH"
vault_command_json "vault write -format=json $PLUGIN/accounts/bob/deploy \
abi=@$ERC20_CONTRACTS_PATH$CONTRACT_SAFE_MATH$ABI_FILE \
bin=@$ERC20_CONTRACTS_PATH$CONTRACT_SAFE_MATH$BIN_FILE"
log_json

header "READ TOKEN $ERC20_ADDRESS BALANCE"
vault_command_json "vault read -format=json $PLUGIN/accounts/bob/erc20/totalSupply contract=$ERC20_ADDRESS"
log_curl_command $(vault read -output-curl-string -format=json $PLUGIN/accounts/bob/erc20/totalSupply contract=$ERC20_ADDRESS)
log_json