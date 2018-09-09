![Immutability](/docs/tagline.png?raw=true "Changes Everything")

Ethereum plugin for Vault
-----------------

The Ethereum secret backend is intended to provide many of the capabilities of an Ethereum wallet. It is designed to support the enterprise adoption of Ethereum though it can be used for a standalone Ethereum wallet. Leveraging [HashiCorp Vault's](https://www.vaultproject.io/) capability set, many different access control patterns are supported. This plugin never exposes the private keys that it manages - all siging operations occur within the Vault secure enclave. The plugin supports smart contract continuous development practices including contract deployment and testing. Some of the functionality (conversion of Ethereum units, retrieving exchange rates, creating accounts and signing transactions) can happen without network connectivity. Other functionality (reading blocks, transactions, account balances and deploying contracts and sending transactions) will require access to an Ethereum RPC interface.

Why is this Important?
-----------------

When an enterprise makes financial transactions, it does so within the context of an institutional apparatus that has many controls to prevent illegitimate loss of funds: if your credit card is stolen, the funds can be returned using the legal system; if wire fraud occurs, you can call the FBI. However, in the world of crypto a loss of your private key(s) can mean a total loss of all funds - there is no regulatory body or instituational apparatus that can change the blockchain. The marriage of Vault and Ethereum allows for fine-grained role based authentication and authorization policies that can provide separation of funds and significant mitigation in the event of a compromise. This mechanism can be leveraged in many ways to build many layers of control. Furthermore the rich features offered by Vault can make auditing, policy definition, and administrative reaction to a breach easy and fun (never thought those words would be in a single sentence).

To put it another way:  If your org is dealing in crypto you need this or something like this.

## No Warrantees Implied

Use of this plugin with real ETH on the mainnet is at your own risk and no warranties should be implied. The guide here describes how to install and run this plugin on a Mac laptop. This plugin can be run on any platform that Vault supports; but, each environment has its own nuances, and for clarity's sake I will only discuss the Mac laptop use case. Running Vault in a production environment in an enterprise requires planning and operational skills. If you would like help running Vault in production, please reach out to [Immutability, LLC](mailto:sales@immutability.io).

## API

[The API is detailed in full here.](./API.md)

Vault is a REST server. Services can be permissioned on granular basis according to their paths. Here is an overview of the paths available and the methods supported:

&nbsp;&nbsp;&nbsp;&nbsp;`└── <MOUNT> `&nbsp;&nbsp;([install](./README.md#install-plugin))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── config `&nbsp;&nbsp;([create](./API.md#create-config), [update](./API.md#update-config), [read](./API.md#read-config))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── accounts `&nbsp;&nbsp;([list](./API.md#list-accounts))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NAME> `&nbsp;&nbsp;([create](./API.md#create-account), [update](./API.md#update-account), [read](./API.md#read-account), [delete](./API.md#delete-account))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── debit `&nbsp;&nbsp;([update](./API.md#debit-account))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── sign `&nbsp;&nbsp;([update](./API.md#sign))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── sign-tx `&nbsp;&nbsp;([update](./API.md#sign-tx))
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── transfer `&nbsp;&nbsp;([update](./API.md#transfer))
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── verify `&nbsp;&nbsp;([update](./API.md#verify))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── addresses `&nbsp;&nbsp;([list](./API.md#list-addresses))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <ADDRESS> `&nbsp;&nbsp;([read](./API.md#read-address))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── balance `&nbsp;&nbsp;([update](./API.md#balance-by-address))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── verify `&nbsp;&nbsp;([update](./API.md#verify-by-address))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── block `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NUMBER> `&nbsp;&nbsp;([read](./API.md#read-block))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── transactions `&nbsp;&nbsp;([read](./API.md#read-block-transactions))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── convert `&nbsp;&nbsp;([update](./API.md#convert))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── export `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NAME> `&nbsp;&nbsp;([create](./API.md#export))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── deploy `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NAME> `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── contracts `&nbsp;&nbsp;([list](./API.md#list-contracts))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │           └── <NAME> `&nbsp;&nbsp;([create](./API.md#deploy-contract), [read](./API.md#read-contract), [delete](./API.md#delete-contract))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── import `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NAME>  `&nbsp;&nbsp;([create](./API.md#import))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── names `&nbsp;&nbsp;([list](./API.md#list-names))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └──  <NAME> `&nbsp;&nbsp;([read](./API.md#read-name))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── balance `&nbsp;&nbsp;([update](./API.md#balance-by-name))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── verify `&nbsp;&nbsp;([update](./API.md#verify-by-name))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    └── transaction `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`        └── <TRANSACTION_HASH> `&nbsp;&nbsp;([read](./API.md#read-transaction))  


## Features

This plugin provides services to:

* Support any Ethereum network (mainnet, Rinkeby, Ropsten, Kovan, private, etc.)
* Create new Ethereum (EOA) accounts.
* Send Ethereum.
* Deploy contracts
* Generate and verify signatures
* Transfer ERC20 tokens.
* Local (to Vault) naming lookups
* Read blocks on the Ethereum blockchain
* Read transactions on Ethereum blockchain
* Create whitelists and blacklists
* Create IP constraints
* Create spending limits (total and per transaction) per account
* Conversion to and from any Ethereum unit (wei, eth, finney, etc.)
* Import and export (Web3 Secret Storage Definition) JSON keystores

Of course, all secrets in Vault are encrypted.

![Vault and Geth Topology](/docs/vault-geth.png?raw=true "Vault Ethereum Plugin")

## Quick start

I recommend that you build the plugin yourself using the standard golang approach:

```
$ go get -u github.com/immutability-io/vault-ethereum
```

This will put the plugin executable in your $GOPATH/bin directory: `$GOPATH/bin/vault-ethereum`. Now you have to install the plugin. For simplicity's sake, I will make some assumptions:

1. You are **NOT** running Vault in `dev` mode. This means that you are running Vault using TLS. [Here an example of how to do this](https://github.com/immutability-io/immutability-project).
2. You have installed Vault as a non-root user.
3. Your Vault configuration and TLS material is at `~/etc/vault.d`. Of course, you can put your Vault configuration anywhere; but, to make these instructions simple.

### Vault configuration

Let's take a look at the Vault configuration directory:

```
$ ls -la1 ~/etc/vault.d
.
..
data
root.crt
vault.crt
vault.hcl
vault.key
vault_plugins
```

And peek at the Vault configuration file:

```
$ cat ~/etc/vault.d/vault.hcl
"default_lease_ttl" = "24h"
"disable_mlock" = "true"
"max_lease_ttl" = "24h"

"backend" "file" {
  "path" = "/Users/cypherhat/etc/vault.d/data"
}

"api_addr" = "https://localhost:8200"
"ui" = "true"
"listener" "tcp" {
  "address" = "localhost:8200"

  "tls_cert_file" = "/Users/cypherhat/etc/vault.d/vault.crt"
  "tls_client_ca_file" = "/Users/cypherhat/etc/vault.d/root.crt"
  "tls_key_file" = "/Users/cypherhat/etc/vault.d/vault.key"
}

"plugin_directory" = "/Users/cypherhat/etc/vault.d/vault_plugins"
```

And make sure our Vault environment variables are set:

```
$ env | grep VAULT
VAULT_ADDR=https://localhost:8200
VAULT_CACERT=/Users/cypherhat/etc/vault.d/root.crt
```

### Install plugin

Now we have enough information to install the plugin. Before we do so, authenticate to Vault as an administrator. Plugins are extremely privileged actors in the Vault ecosystem, so they need to be installed and confured by the Vault admin.

We will move the vault-ethereum plugin to the `plugin_directory`, add the plugin to Vault's plugin catalog and then mount and enable it. Note that we sign the plugin executable.:

```sh
$ mv $GOPATH/bin/vault-ethereum $HOME/etc/vault.d/vault_plugins
$ export SHASUM256=$(shasum -a 256 "$HOME/etc/vault.d/vault_plugins/vault-ethereum" | cut -d' ' -f1)
$ vault write sys/plugins/catalog/ethereum-plugin \
      sha_256="${SHASUM256}" \
      command="vault-ethereum --ca-cert=$HOME/etc/vault.d/root.crt --client-cert=$HOME/etc/vault.d/vault.crt --client-key=$HOME/etc/vault.d/vault.key"
$ vault secrets enable -path=ethereum -description="Immutability's Ethereum Wallet" -plugin-name=ethereum-plugin plugin
```

We should now be able to see the plugin if we query Vault for the available secrets backends (there are several listed here - the one of interest is `ethereum`:

```
$ vault secrets list
Path            Type         Accessor              Description
----            ----         --------              -----------
aws/            aws          aws_2e82318b          AWS credentials
btc/mainnet/    plugin       plugin_d0630fca       BTC Mainnet Wallet
btc/regtest/    plugin       plugin_e18ad345       BTC Regression Test Wallet
btc/simnet/     plugin       plugin_9e9c7181       BTC Simnet Wallet
btc/testnet/    plugin       plugin_b735aa97       BTC Testnet Wallet
cubbyhole/      cubbyhole    cubbyhole_db313a57    per-token private secret storage
ethereum/       plugin       plugin_69372d75       Immutability's Ethereum Wallet
identity/       identity     identity_4638066d     identity store
ltc/            plugin       plugin_3f43f7c3       LTC Wallet
secret/         kv           kv_44746ed8           key/value secret storage
sys/            system       system_2a5f140a       system endpoints used for control, policy and debugging
trust/          plugin       plugin_0cf966e2       Immutability's Trustee Service
```

## Playing with Immutability's Ethereum Wallet

Before we do anything with the plugin, we need to configure it: we need to tell the plugin which Ethereum network it will use (Rinkeby by default) and what node address to use for RPC communication (`https://rinkeby.infura.io` by default.) Optionally, we can add the [Infura API key](https://infura.io/register). Since we are running Vault in a trusted environment, we may (optionally) restrict network access by client IP address. [See the API for more details](./API.md#create-config).

### Configure the plugin

If we want to accept all defaults, we can configure the plugin as follows:

```sh
$ vault write -f ethereum/config
Key                Value
---                -----
api_key            n/a
bound_cidr_list    <nil>
chain_id           4
rpc_url            https://rinkeby.infura.io
```

Now we can start to play! 

## Unauthenticated endpoints

Probably the most controversial aspect of the design of this plugin is the fact that we have a handful of unauthenticated endpoints. Since Vault is a tool for managing secrets, why would we want it to allow unauthenticated access to anything? The answer is that we want the experience of interacting with the Ethereum ecosystem with Vault to be easy (and fun.) Having a consistent API for both trusted and untrusted actions increases usability. I believe the Vault plugin model allows Vault to be a **platform** for blockchain applications and development.

Since Immutability's Ethereum Wallet has several unauthenticated endpoints ([detailed in the API](./API.md)), we will play with a few here before we create any Ethereum accounts:

### ETH Unit Converter

[There is a website](https://etherconverter.online/) that will convert any ETH unit to any other. Since Immutability's Ethereum Wallet only allows you to send ETH in wei, I thought it would be useful to replicate this capability. **NOTE** we can also get the current average exchange value but converting any unit to USD. We can also convert from USD to any unit:

```
$ vault write ethereum/convert unit_from="wei" unit_to="tether" amount="4.4"
Key            Value
---            -----
amount_from    4.4
amount_to      0.0000000000000000000000000000044
unit_from      wei
unit_to        tera

$ vault write ethereum/convert unit_from="eth" unit_to="wei" amount="4.4"
Key            Value
---            -----
amount_from    4.4
amount_to      4400000000000000000
unit_from      ether
unit_to        wei

$ vault write ethereum/convert amount=1 unit_from=eth unit_to=usd
Key            Value
---            -----
amount_from    1
amount_to      224.982648392
unit_from      ether
unit_to        usd
```

All known ETH units are supported.

### Query Block By Number

Sometimes it is useful to know whether a block exists or what its hash was or other details:

```
$ vault read -format=json ethereum/block/24395834800568
No value found at ethereum/block/24395834800568

$ vault read -format=json ethereum/block/2800569
{
  "request_id": "c02c08a7-3283-8b17-a486-ddab3321bb16",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "block": 2800569,
    "block_hash": "0xf403508f16c234a7603a04e968409ceca8db4c36425edefa3bcee2215f06b8fd",
    "difficulty": 1,
    "time": 1534078643,
    "transaction_count": 3
  },
  "warnings": null
}
```

### Read Transactions By Block Number

If you want to know the transaction hashes at a particular block, you append `transactions` to the path above:

```
$ vault read -format=json ethereum/block/2800569/transactions
{
  "request_id": "ba77bd3b-923b-1ad6-a289-d80f9752934b",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "0x577440e22b85895fb852a326d83e5fc9720f7f6e6037196806e82cd7a1e428cb": {
      "address_to": "0xd35CB9dFf4b94d7Ab9bc0A8f221b127a09429AA7",
      "gas": 200000,
      "gas_price": 1000000000,
      "nonce": 88582,
      "value": "0"
    },
    "0x84442dc3f585436fc34d5881988c3612a4d7782cb50e1e54f2d06082b5b307e3": {
      "address_to": "0xa9Ad608533A7817456c706286202405DcEF471b3",
      "gas": 52744,
      "gas_price": 1000000000,
      "nonce": 195,
      "value": "0"
    },
    "0xbbc49ad70961b891949964ebe14caa53434cab2f3c66cf5f8275a4b06f5e5d1a": {
      "address_to": "0x64A10C83d5EF68301625F4df0fACB38C78d622E4",
      "gas": 6500000,
      "gas_price": 9000000000,
      "nonce": 962,
      "value": "0"
    }
  },
  "warnings": null
}
```

### Read Transaction Details

This listing of transactions by block omits certain details, so there is another method to return more details about the transaction (or whether a transaction exists):

```
$ vault read -format=json ethereum/transaction/0xbbc49ad70961b891949964ebe14caa53434cab2f3c66cf5f8275a4b06f5e5d1a
{
  "request_id": "6fa9bca9-c47c-41ab-b8b4-adf750f48e20",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "address_from": "0xAc09AdC0939754C2D4D16D849d9C3D7526b7f9f2",
    "address_to": "0x64A10C83d5EF68301625F4df0fACB38C78d622E4",
    "gas": 6500000,
    "gas_price": 9000000000,
    "nonce": 962,
    "pending": false,
    "receipt_status": 1,
    "transaction_hash": "0xbbc49ad70961b891949964ebe14caa53434cab2f3c66cf5f8275a4b06f5e5d1a",
    "value": "0"
  },
  "warnings": null
}

$ vault read -format=json ethereum/transaction/0xbbc49ad70961b891949964ebe14caa53434cab2f3c66cf5f8275a4b06f5e51a
No value found at ethereum/transaction/0xbbc49ad70961b891949964ebe14caa53434cab2f3c66cf5f8275a4b06f5e51a
```

## Ethereum Accounts

In order to create and use Ethereum accounts, you have to be authenticated. The power of Vault lies in its ability to broker many different authentication mechanisms to allow practical authorization checks for any path. I will not go into detail here on policy design. I simply assume that the Vault client is authenticated and permitted to perform the following actions. If you would like guidance on how to design a Vault permission model for your enterprise, please reach out to [Immutability, LLC](mailto:sales@immutability.io).  

### Create Account

Let's create an Ethereum Account. It is very simple if we want to accept the defaults. 

```
$ vault write -f ethereum/accounts/muchwow
Key                     Value
---                     -----
address                 0x7b715f8748ef586b98d3e7c88f326b5a8f409cd8
blacklist               <nil>
spending_limit_total    0
spending_limit_tx       0
total_spend             0
whitelist               <nil>
```

For more details on whitelist, blacklist, and spending limits please [refer to the API documentation](./API.md#create-account).

Notice that the private key is not returned. Duh. The only way to exfiltrate the private key from Vault is to use the [export](./API.md#export) feature. 

### Read Account

Let's take a look at the account we just created:

```
$ vault read ethereum/accounts/muchwow
Key                     Value
---                     -----
address                 0x7b715f8748ef586b98d3e7c88f326b5a8f409cd8
balance                 1000000000000000000
balance_in_usd          0
blacklist               <nil>
spending_limit_total    0
spending_limit_tx       0
total_spend             0
whitelist               <nil>
```

Wait! How in the heck did that balance get there? Since wei is used across every service in this plugin, we can use our handy dandy conversion service to see what that is in ETH. (I know, if you can't do this math, maybe crypto isn't for you... says to mirror everyday...). Note also that the `balance_in_usd` is `0`. This is because we are on a test net (Rinkeby) where the ETH has no exchange value. If this had been an account read from the mainnet, then an attempt at estimating the value of the ETH in USD would have been made:

```
$ vault write ethereum/convert unit_from="wei" unit_to="eth" amount="1000000000000000000"
Key            Value
---            -----
amount_from    1000000000000000000
amount_to      1
unit_from      wei
unit_to        ether
```

So, somebody sent us 1 ETH? How did that happen? Well, it wasn't magic and it wasn't the donation of a generous crypto billionaire (sadly). And, since this is Rinkeby, it is worthless. I used the [Rinkeby faucet](https://faucet.rinkeby.io/).

### Send ETH

Well, since the `muchwow` account is flush with funds, let's pay it forward to a new account. First, we will create an account to move ETH to:

```
$ vault write -f ethereum/accounts/lesswow
Key                     Value
---                     -----
address                 0x36d1f896e55a6577c62fdd6b84fbf74582266700
blacklist               <nil>
spending_limit_total    0
spending_limit_tx       0
total_spend             0
whitelist               <nil>
```

I want to send 0.2 ETH from `muchwow` to `lesswow`. Just because I don't trust myself to type properly, I use our handy dandy conversion service to do this:

```$ vault write ethereum/convert unit_from="eth" unit_to="wei" amount="0.2"
Key            Value
---            -----
amount_from    0.2
amount_to      200000000000000000
unit_from      ether
unit_to        wei
```

Now, we send the ETH from the `muchwow` account to the address of the `lesswow` account. I use suggested gas price and the default gas limit of 21000.

```

$ vault write ethereum/accounts/muchwow/debit amount=200000000000000000 address_to="0x36d1f896e55a6577c62fdd6b84fbf74582266700"
Key                       Value
---                       -----
amount                    200000000000000000
amount_in_usd             0
address_from              0x7b715f8748ef586b98d3e7c88f326b5a8f409cd8
address_to                0x36D1F896E55a6577C62FDD6b84fbF74582266700
gas_limit                 21000
gas_price                 2000000000
signed_transaction        0xf86b07843b9aca00825208948440a3f9243b96cd934de1b7a400368d880b041d88016345785d8a0000802ca0c8f2511d337ce9180deb525fb714100157a89fba7e990677f74609dde74bad21a0239b33f6439b28a9f9d9e8cbbb45fd77e1c1dddc33d16c07345667ee12d2a767
starting_balance          1000000000000000000
starting_balance_in_usd   0
total_spend               200000000000000000
transaction_hash          0x0b4938a1a44f545deeea500d50761c22bfe2bc006b26be8adf4dcd4fc0597769
```

**NOTE**: The `starting_balance` is the balance **before** the transaction is committed. Let's check 2 things - the transaction and the account balance of `lesswow`.

Read the transaction details:

```
$ vault read -format=json ethereum/transaction/0x0b4938a1a44f545deeea500d50761c22bfe2bc006b26be8adf4dcd4fc0597769
{
  "request_id": "960eab1a-4a5a-cbdd-1ebb-716a5bf5c872",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "address_from": "0x7B715f8748Ef586b98d3e7c88f326b5a8F409Cd8",
    "address_to": "0x36D1F896E55a6577C62FDD6b84fbF74582266700",
    "gas": 21000,
    "gas_price": 2000000000,
    "nonce": 0,
    "pending": false,
    "receipt_status": 1,
    "transaction_hash": "0x0b4938a1a44f545deeea500d50761c22bfe2bc006b26be8adf4dcd4fc0597769",
    "value": "200000000000000000"
  },
  "warnings": null
}
```

Read the balance:

```
$ vault read ethereum/accounts/lesswow
Key                     Value
---                     -----
address                 0x36d1f896e55a6577c62fdd6b84fbf74582266700
balance                 200000000000000000
balance_in_usd          0
blacklist               <nil>
spending_limit_total    0
spending_limit_tx       0
total_spend             0
whitelist               <nil>
```

As you can see, `lesswow` is more wow than previously. In fact, it seems wrong that `lesswow` is named as such. Let's *rename* it to `morewow`. This will involve exporting and importing.

### Export 

Let's export the private key associated with `lesswow` to a (Web3 Secret Storage Definition) JSON keystore. Truth be told, this was the most annoying part of developing this plugin - the [go-ethereum codebase](https://github.com/ethereum/go-ethereum) mysteriously chose not to export certain classes, so there was a bit of cutting and carving to get the ability to render a (Web3 Secret Storage Definition) JSON keystore without getting tangled up with their wallet implementations.

Let's export the `lesswow` private key:

```
$ vault write ethereum/export/lesswow path=$(pwd)
Key           Value
---           -----
passphrase    vibes-subsystem-print-truck-nuptials-outcome-drudge-setting-raging
path          /Users/cypherhat/develop/go/src/github.com/immutability-io/vault-ethereum/test/ethereum/config/UTC--2018-08-19T17-04-42.950941035Z--36d1f896e55a6577c62fdd6b84fbf74582266700
```

The first thing to notice here is that the passphrase used to encrypt this keystore is generated for you. The second thing to notice is that this passphrase is displayed for all to see. In a real use case, you would use your considerable bash skills to prevent this exposure.

There is a [nifty Golang utility](https://github.com/atotto/clipboard) that works in OSX, Windows and Linux that allows you to pipe the output of a command into the clipboard. I have installed this (as well as the [awesome jq utility](https://github.com/stedolan/jq)) in my environment. This allows me to do this:

```sh
$ vault write -format=json ethereum/export/lesswow path=$(pwd) | jq -r .data.passphrase | gocopy
```

When you export like this, nothing is revealed to the screen but the JSON keystore is exported and your clipboard contains the passphrase. Then it is a simple matter to import the JSON keystore into something like MetaMask. Each time a private key is exported, it is encrypted with a newly generated passphrase.

### Import

Import is the reverse of export. It will create a new named account in Vault using the same private key. 

```
$ vault write ethereum/import/morewow path="/Users/cypherhat/develop/go/src/github.com/immutability-io/vault-ethereum/test/ethereum/config/UTC--2018-08-19T17-04-42.950941035Z--36d1f896e55a6577c62fdd6b84fbf74582266700" passphrase=$PASSPHRASE
Key                     Value
---                     -----
address                 0x36d1f896e55a6577c62fdd6b84fbf74582266700
blacklist               <nil>
spending_limit_total    n/a
spending_limit_tx       n/a
total_spend             n/a
whitelist               <nil>
```

If we read from this account, we will see that the balance is indeed morewow:

```
$ vault read ethereum/accounts/morewow
Key                     Value
---                     -----
address                 0x36d1f896e55a6577c62fdd6b84fbf74582266700
balance                 200000000000000000
balance_in_usd          0
blacklist               <nil>
spending_limit_total    n/a
spending_limit_tx       n/a
total_spend             n/a
whitelist               <nil>
```

### Cross References

We just created 2 accounts - 2 of which have the same private key (Ethereum account) but different names. We can see this by listing the names of the accounts we have created:

```
$ vault list ethereum/names
Keys
----
lesswow
morewow
muchwow
```

We can also see the actual accounts:

```
$ vault list ethereum/addresses
Keys
----
0x36d1f896e55a6577c62fdd6b84fbf74582266700
0x7b715f8748ef586b98d3e7c88f326b5a8f409cd8
```

And to see the addresses by name:

```
$ vault read ethereum/names/muchwow
Key        Value
---        -----
address    0x7b715f8748ef586b98d3e7c88f326b5a8f409cd8

$ vault read ethereum/names/morewow
Key        Value
---        -----
address    0x36d1f896e55a6577c62fdd6b84fbf74582266700

$ vault read ethereum/names/lesswow
Key        Value
---        -----
address    0x36d1f896e55a6577c62fdd6b84fbf74582266700
```

And, of course we can do the opposite:

```
 $ vault read ethereum/addresses/0x36d1f896e55a6577c62fdd6b84fbf74582266700
Key      Value
---      -----
names    [morewow lesswow]

$ vault read ethereum/addresses/0x7b715f8748ef586b98d3e7c88f326b5a8f409cd8
Key      Value
---      -----
names    [muchwow]
```

## Create contracts

Suppose you have written a smart contract. Likely, it is only one or 2 deployment cycles away from yielding ICO riches. So, you better deploy it. The Vault plugin allows you to deploy a compiled smart contract.

Sending any transaction on the Ethereum network requires the payment of fees. So, you send the transaction that deploys a contract **from** an Ethereum account with a positive balance.

Assume that the compiled contract is the file `Helloworld.bin`. Deployment is simple. We have to decide how much to fund our contract with - 0.1 ETH... and how much is that in wei? And we have to provide a gas limit (though one will be suggested if not supplied):

```sh

$ vault write ethereum/convert unit_from="wei" unit_to="eth" amount="100000000000000000"
Key            Value
---            -----
amount_from    100000000000000000
amount_to      0.1
unit_from      wei
unit_to        ether

$ vault write ethereum/accounts/muchwow/contracts/helloworld transaction_data=@Helloworld.bin amount=100000000000000000 gas_limit=1500000
```

The above command says: *Deploy a contract, named `helloworld`, from the account named `muchwow` and fund it with 0.1 ETH and provide a gas limit of 1500000 wei*

```
Key                 Value
---                 -----
transaction_hash    0x2ff5dd013e5a4d00cf007a7fb689c4ebf50541c2e7ddfaf16212e7ed1ba70f4c

```

When you deploy a contract, the contract address isn't immediately available. What is returned from the Vault-Ethereum plugin after a contract deployment is just:

* `transaction_hash`: The hash of the contract deployment transaction.

### Read contract address

Since the contract address isn't known at the point when the transaction is sent, so you have to **revisit** the contract (with a read operation) to determine the address:

```sh
$ vault read ethereum/accounts/muchwow/contracts/helloworld
```

```
Key                 Value
---                 -----
address             0x78545F1100912B001418741177b5b1eFB00DfaF1
transaction_hash    0x2ff5dd013e5a4d00cf007a7fb689c4ebf50541c2e7ddfaf16212e7ed1ba70f4c
```

## More Use Cases

There are many more uses cases to be explored; but, I will leave that to you. If you brave enough...

## Credits

None of this would have been possible without the fantastic [tutorial](https://www.hashicorp.com/blog/building-a-vault-secure-plugin) on Vault Plugins by Seth Vargo. Seth is one of those rare individuals who can communicate the simple essence of a complex technology in practical terms.

I had the great fortune to attend DevCon3 in November, 2017 and hear Andy Milenius speak with clarity and vision about how the Ethereum developer ecosystem should embrace the Unix philosophy - the same philosophy that makes **everything-as-code** possible: simple tools, with clear focus and purpose, driven by repeatable and interoperable mechanics. So, when I returned from DevCon3 (and dug out from my work backlog - a week away is hard) I installed `seth` and `dapp` and found inspiration.

The community chat that the [dapphub](https://dapphub.com/) guys run (esp. Andy and Mikael and Daniel Brockman) is a super warm and welcoming place that pointed me towards code that greatly helped this experiment.

Last but not least, I have to thank Miguel Mota for his wonderfully elegant [Ethereum Development with Go produced by Miguel Mota](https://github.com/miguelmota/ethereum-development-with-go-book).

## License

This code is licensed under the Apache 2 license. Please feel free to use it. Please feel free to contribute.

## Donations?

Send ETH to 0x4169c9508728285e8A9f7945D08645Bb6b3576e5 and you will be blessed in the next life.

![Donations Accepted](/docs/0x4169c9508728285e8A9f7945D08645Bb6b3576e5.png?raw=true "0x4169c9508728285e8A9f7945D08645Bb6b3576e5")

