#!/bin/bash

MNEMONIC="volcano story trust file before member board recycle always draw fiction when"
CHAIN_ID=5777
ROPSTEN_CHAIN_ID=3
PORT=8545
RPC_URL="http://ganache:$PORT"
ROPSTEN_URL="https://ropsten.infura.io/v3/765a9b4697d8438ba1e73c8c2d0838e6"
ERC721_ADDRESS="0xd8bbF8cEb445De814Fb47547436b3CFeecaDD4ec"
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

function header() {
  printf "### %s:\n" "$@"
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
  JSON=$(eval "$@" | jq .data)
  #echo $JSON | jq .
}

header "HOW TO CONFIGURE THE PLUGIN AFTER IT HAS BEEN ENABLED"
vault_command_json "vault write -format=json $PLUGIN/config  rpc_url=$RPC_URL chain_id=$CHAIN_ID"
log_json

header "HOW TO CREATE AN ACCOUNT NAMED BOB USING A MNEMONIC"
vault_command_json "vault write -format=json '$PLUGIN'/accounts/bob mnemonic='$MNEMONIC'"
log_json

header "1: GET  BOBS ADDRESS"
vault_command_json "vault read -format=json vault-ethereum/accounts/bob"
BOB_ADDRESS=$(echo $JSON | jq .address)
log_value $BOB_ADDRESS

header "HOW TO CREATE AN ACCOUNT NAMED ALICE WITH NO MNEMONIC"
vault_command_json "vault write -f -format=json $PLUGIN/accounts/alice"
log_json

header "HOW TO TRANSFER 0.5 ETH FROM BOB TO ALICE"
header "1: CONVERT 0.5 ETH TO WEI"

vault_command_json "vault write -format=json vault-ethereum/convert unit_from=ETH unit_to=WEI amount=0.5"
log_json
AMOUNT_TO=$(echo $JSON | jq .amount_to)
echo "Amount $AMOUNT_TO"

header "2: GET ALICES ADDRESS"
vault_command_json "vault read -format=json vault-ethereum/accounts/alice"
ALICE_ADDRESS=$(echo $JSON | jq .address)
log_value $ALICE_ADDRESS

header "3: SEND BOBS ETH TO ALICE"
vault_command_json "vault write -format=json vault-ethereum/accounts/bob/transfer to=$ALICE_ADDRESS amount=$AMOUNT_TO"
log_json

header "HOW TO CHECK BALANCE OF THE ACCOUNT NAMED BOB"
vault_command_json "vault read -format=json '$PLUGIN'/accounts/bob/balance"
log_json

header "HOW TO SIGN A TRANSACTION WITH DATA"
vault_command_json "vault write -format=json '$PLUGIN'/accounts/alice/sign-tx to='$BOB_ADDRESS' data='$RAW_TX' amount=20000000000000000"
log_json

header "HOW TO SIGN MESSAGE. Message signature can be verify on: https://etherscan.io/verifySig/2156" 
vault_command_json "vault write -format=json '$PLUGIN'/accounts/bob/sign message='$MESSAGE'"
log_json

header "DEPLOY CONTRACT $CONTRACT_FIXED_SUPPLY_TOKEN"
vault_command_json "vault write -format=json '$PLUGIN'/accounts/bob/deploy \
abi=@'$ERC20_CONTRACTS_PATH$CONTRACT_FIXED_SUPPLY_TOKEN$ABI_FILE' \
bin=@'$ERC20_CONTRACTS_PATH$CONTRACT_FIXED_SUPPLY_TOKEN$BIN_FILE'"
log_json
ERC20_ADDRESS=$(echo $JSON | jq -r .contract)

header "DEPLOY CONTRACT $CONTRACT_OWNED"
vault_command_json "vault write -format=json '$PLUGIN'/accounts/bob/deploy \
abi=@'$ERC20_CONTRACTS_PATH$CONTRACT_OWNED$ABI_FILE' \
bin=@'$ERC20_CONTRACTS_PATH$CONTRACT_OWNED$BIN_FILE'"
log_json

header "DEPLOY CONTRACT $CONTRACT_SAFE_MATH"
vault_command_json "vault write -format=json '$PLUGIN'/accounts/bob/deploy \
abi=@'$ERC20_CONTRACTS_PATH$CONTRACT_SAFE_MATH$ABI_FILE' \
bin=@'$ERC20_CONTRACTS_PATH$CONTRACT_SAFE_MATH$BIN_FILE'"
log_json

header "READ TOKEN $ERC20_ADDRESS TOTAL SUPPLY"
vault_command_json "vault read -format=json '$PLUGIN'/accounts/bob/erc20/totalSupply contract='$ERC20_ADDRESS'"
log_json

header "TEST ON ROPSTEN TESTNET"

header "RE-CONFIGURE THE PLUGIN (ROPSTEN TESTNET)"
vault_command_json "vault write -format=json '$PLUGIN'/config  rpc_url='$ROPSTEN_URL' chain_id='$ROPSTEN_CHAIN_ID'"
log_json

header "READ TOKEN $ERC721_ADDRESS TOTAL SUPPLY"
vault_command_json "vault read -format=json $PLUGIN/accounts/bob/erc721/totalSupply contract='$ERC721_ADDRESS'"
log_json
