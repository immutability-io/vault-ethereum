Ethereum plugin for Vault
-----------------

The Ethereum secret backend is intended to provide many of the capabilities of an Ethereum wallet. It is designed to support smart contract continuous development practices including contract deployment and testing. It has only been exercised on private Ethereum chains and the Rinkeby testnet. Some of the functionality (creating accounts and signing contract creation transactions) can happen without any local `geth` node to be running. Other functionality (deploying contracts and sending transactions - still in development) will require the geth RPC interface.

## Works for me, but...

This plugin is still in the early stages of development. I have used it extensively on private chains and on Rinkeby. I am HODLing the real ETH I have, so I have only used this plugin to check the balance of my mainnet accounts. Use of this plugin with real ETH on the mainnet is at your own risk and no warranties should be implied. 

## Features

This plugin provides services to:

* Create new externally controlled accounts.
* Import JSON keystores (with provided passphrase.)
* Export JSON keystores
* Sign transactions for contract deployment
* Sign arbitrary data
* Send Ethereum
* Deploy contracts

All secrets in Vault are encrypted. However, for ease of integration with `geth`, the plugin stores the Ethereum private key in encrypted (JSON keystore) format.

![Vault and Geth Topology](/doc/vault-geth.png?raw=true "Vault Ethereum Plugin")

## Quick start

Building the plugin and installing it will be covered later, but let's assume that has been done. It is important to note that the Vault Ethereum plugin can be mounted at any path. A common model is to use a well defined namespace for mounting Vault backends - for example, using the GitHub org/repo namespace: `ethereum/immutability-io/world-domination-token`. For this discussion, we assume that the Vault Ethereum plugin has been mounted at `ethereum`.

### Create new externally controlled accounts

Let's create an Ethereum account using our private chain (w/ chain ID `1977`):

```sh
$ vault write ethereum/accounts/test chain_id=1977
```

That's all that is needed. **NOTE:** we also assume that the vault client has been authenticated and has permission `write` to `ethereum/accounts/test`. Discussion of the [Vault authentication model can be found here](https://www.vaultproject.io/docs/concepts/auth.html).

The Ethereum plugin will return the following information from the above command:

```
Key          Value
---          -----
address      0x98FE39097a0eA57A2D8bC4f7A654F180b962be7f
blacklist    <nil>
chain_id     1977
rpc_url      http://localhost:8545
whitelist    <nil>
```

The parameters `rpc_url` defaults to `http://localhost:8545` but it can be customized when creating an account.

### Read externally controlled accounts

We can read the account stored at `ethereum/accounts/test` as follows:

```sh
$ vault read ethereum/accounts/test
```

```
Key             	Value
---             	-----
address      0x98FE39097a0eA57A2D8bC4f7A654F180b962be7f
blacklist    <nil>
chain_id     1977
rpc_url      http://localhost:8545
whitelist    <nil>
```

For convenience, we can read the balance for this account:

```
$ vault read ethereum/accounts/test/balance
Key                 Value
---                 -----
address             0x98FE39097a0eA57A2D8bC4f7A654F180b962be7f
pending_balance     10000000000000000000
pending_nonce       0
pending_tx_count    0
```

### Create contracts

Suppose you have written a smart contract. Likely, it is only one or 2 deployment cycles away from yielding ICO riches. So, you better deploy it. The Vault plugin allows you to deploy a compiled smart contract.

Sending any transaction on the Ethereum network requires the payment of fees. So, you send the transaction that deploys a contract **from** an Ethereum account with a positive balance.

Assume that the compiled contract is the file `Helloworld.bin`. Deployment is simple:

```sh
$ vault write ethereum/accounts/test/contracts/helloworld transaction_data=@Helloworld.bin value=10000000000000000000 gas_price=21000000000 gas_limit=1500000
```

The above command says: *Deploy a contract, named `helloworld`, from the account named `test`*

```
Key        Value
---        -----
tx_hash    0x2ff5dd013e5a4d00cf007a7fb689c4ebf50541c2e7ddfaf16212e7ed1ba70f4c

```

When you deploy a contract, the contract address isn't immediately available. What is returned from the Vault-Ethereum plugin after a contract deployment is just:

* `tx_hash`: The hash of the contract deployment transaction.

### Read contract address

Since the contract address isn't known at the point when the transaction is sent, so you have to **revisit** the contract (with a read operation) to determine the address:

```sh
$ vault read ethereum/accounts/test/contracts/helloworld
```

```
Key        Value
---        -----
address    0x78545F1100912B001418741177b5b1eFB00DfaF1
tx_hash    0x2ff5dd013e5a4d00cf007a7fb689c4ebf50541c2e7ddfaf16212e7ed1ba70f4c
```

### Import keystores from other wallets

Lastly, suppose we already have an Ethereum wallet and we are convinced that storing the private key in Vault is something we want to do. The plugin supports importing JSON keystores. **NOTE:** you have to provide the path to a single keystore - this plugin doesn't support importing an entire directory yet.

```sh
$ ls -la ~/.ethereum/keystore
```

```
total 24
drwxr-xr-x  5 immutability  admin  170 Dec  2 11:57 .
drwxr-xr-x  3 immutability  admin  102 Dec  2 11:55 ..
-rw-r--r--  1 immutability  admin  492 Dec  2 11:56 UTC--2017-12-01T23-13-37.315592353Z--a152e7a09267bcff6c33388caab403b76b889939
-rw-r--r--  1 immutability  admin  492 Dec  2 11:56 UTC--2017-12-01T23-13-56.838050955Z--0374e76da2f0be85a9fdc6763864c1087e6ed28b
-rw-r--r--  1 immutability  admin  492 Dec  2 11:57 UTC--2017-12-01T23-14-16.032409548Z--f19a9a9b2ad60c66429451075046869a9b7014f7
```

As will be discussed in the next section, handling passphrases is always problematic. Care should be taken when importing a keystore not to leak the passphrase to the shell's history file or to the environment:

```sh
$ read -s PASSPHRASE; vault write ethereum/import/test2 path=/Users/immutability/.ethereum/keystore/UTC--2017-12-01T23-13-37.315592353Z--a152e7a09267bcff6c33388caab403b76b889939 passphrase=$PASSPHRASE; unset PASSPHRASE
```

```
Key    	Value
---    	-----
address	0xa152E7a09267bcFf6C33388cAab403b76B889939
```

Now, we can use the imported account as we did with our generated account.

### Export JSON keystores

If we wish to export the account (private key) as a JSON keystore if we want to import it into a external wallet. Each time we export it, a cryptographically strong passphrase is generated to encrypt the JSON keystore. **You need to take care not to reveal this passphrase**:

```sh
$ vault write ethereum/accounts/test/export path=$(pwd)
```

```
Key 	         Value
--- 	         -----
passphrase     resource-ladybug-subzero-childish-nutrient-flyer-macaw-whacky-flagship
path	         /Users/immutability/.ethereum/keystore/UTC--2017-12-01T23-13-37.315592353Z--a152e7a09267bcff6c33388caab403b76b889939
```

#### Export passphrase to clipboard

There is a [nifty Golang utility](https://github.com/atotto/clipboard) that works in OSX, Windows and Linux that allows you to pipe the output of a command into the clipboard. I have installed this (as well as the [awesome jq utility](https://github.com/stedolan/jq)) in my environment. This allows me to do this:

```sh
$ vault write -format=json ethereum/accounts/test/export path=$(pwd) | jq .data.passphrase | tr -d '"' | gocopy
```

When you export like this, nothing is revealed to the screen but the JSON keystore is exported and your clipboard contains the passphrase. Then it is a simple matter to import the JSON keystore into something like MetaMask:

![Import into MetaMask](/doc/import_metamask.png?raw=true "How to use an Exported JSON Keystor")

### Signing arbitrary data

We can also sign arbitrary data using the `sign` endpoint:

```sh
$ vault write ethereum/accounts/test2/sign  data=@../data/test.txt
```

```
Key      	Value
---      	-----
signature	0xe81d649f2a295aa58ad0d67b2adf0f5f336e11a46bd69347f197f073244863406027daed083675b5af5c99b3f1608b53620cd02ca51a65b67773b1580552deb501
```

### Verifying signatures

Of course, we can also verify signatures using the `verify` endpoint. This will verify that a particular account (in this case the account named `test`) actually signed the data:

```sh
$ vault write ethereum/accounts/test/verify data=@somefile.txt signature=0xdb6f22f068ae23473beb9b71bb1a2df64a71cb2e51fc43d67558ba8934da572d49b3faa8da387703870474c92beb8c53e89bbd02ba2356b5fc8fa5b342d8fb7b00
Key         Value
---         -----
verified    true
```

### Sending Ethereum

Now that we have accounts in Vault, we can drain them! We can send ETH to other accounts on the network. (In my test case, it must be emphasized: this is a private test network or Rinkeby.) Assuming there are funds in the account, we can send ETH to another address. In this case, we write a `debit` to the `test3` account:

```sh
$ vault write ethereum/accounts/test3/debit to=0x0374E76DA2f0bE85a9FdC6763864c1087e6Ed28b value=10000000000000000000
```

```
Key    	Value
---    	-----
tx_hash	0xe99f3de1dfbae82121a009b9d3a2a60174f2904721ec114a8fc5454a96e62ba8

```

If the gas limit is omitted, we will try to estimate it; if the gas price is omitted, we will use a suggested gas price.

#### Rudimentary controls

I have implemented a few rudimentary controls to prevent sending transactions that shouldn't be sent.

##### Insufficient funds

If your pending balance at the point in time you are trying to send ETH is less than the amount of ETH you want to send, then the transaction will not be attempted.

```sh
$ vault write ethereum/accounts/etherbase/debit to=0xD9E025bFb6ef48919D9C1a49834b7BA859714cD8 amount=1000000000000000000000000
Error writing data to ethereum/accounts/etherbase/debit: Error making API request.

URL: PUT https://localhost:8200/v1/ethereum/accounts/etherbase/debit
Code: 500. Errors:

* 1 error occurred:

* Insufficient funds to debit 1000000000000000000000000 because the current account balance is 244998580000000000000
```

##### Whitelisting accounts

Imagine that there is an approval process that determines the ETH addresses that you are allowed to send ETH to. This is implemented in the plugin by setting a `whitelist` attribute on the Ethereum account:

```sh
$ vault write ethereum/accounts/etherbase whitelist=0xD9E025bFb6ef48919D9C1a49834b7BA859714cD8,0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572
Key          Value
---          -----
address      0x3943FF61FF803316cF02938b5b0b3Ba3bbE183e4
blacklist    <nil>
chain_id     4
rpc_url      http://localhost:8545
whitelist    [0xD9E025bFb6ef48919D9C1a49834b7BA859714cD8 0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572]
```

If you set a whitelist, you will be prevented from sending funds to an account not in the whitelist:

```sh
$ vault write ethereum/accounts/etherbase/debit to=0x16B429f9B46Bc50B375660a4aFe7e07b6369D8aC amount=100000000000000000
Error writing data to ethereum/accounts/etherbase/debit: Error making API request.

URL: PUT https://localhost:8200/v1/ethereum/accounts/etherbase/debit
Code: 500. Errors:

* 1 error occurred:

* 0x16B429f9B46Bc50B375660a4aFe7e07b6369D8aC is not in the whitelist
```

##### Blacklisting accounts

Imagine that there is list of accounts that are known to authorities as being vehicles for money laundering. (Hint: this is an excellent business idea.) If you want to prevent any ETH from being sent to these bad accounts you can create a blacklist:

```sh
$ vault write ethereum/accounts/etherbase blacklist=0xD9E025bFb6ef48919D9C1a49834b7BA859714cD8,0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572
Key          Value
---          -----
address      0x3943FF61FF803316cF02938b5b0b3Ba3bbE183e4
blacklist    [0xD9E025bFb6ef48919D9C1a49834b7BA859714cD8 0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572]
chain_id     4
rpc_url      http://localhost:8545
whitelist    <nil>
```

If you set a blacklist, you will be prevented from sending funds to an account in the blacklist:

```sh
$ vault write ethereum/accounts/etherbase/debit to=0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572 amount=100000000000000000
Error writing data to ethereum/accounts/etherbase/debit: Error making API request.

URL: PUT https://localhost:8200/v1/ethereum/accounts/etherbase/debit
Code: 500. Errors:

* 1 error occurred:

* 0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572 is blacklisted
```

## Passphrases and authentication

In the Cryptocurrency universe, private keys are protected by passphrases. If you forget your passphrase, you lose access to your private keys and all your funds. Authentication via passphrase is often very awkward; however, it is a choice. The beauty of Vault is that it lets you make that choice. If you want to use passphrases to protect your private keys, you can: simply set up a `userpass` authentication backend. Then, if you so desire, you can add MFA to that authentication mechanism (try that with Mist or MetaMask.) This also ignores the fact that the Vault unsealing process is yet another control that prevents unauthorized use of private keys.

The net result is: if you want to use Vault as a personal Ethereum wallet, it has better controls than most software alternatives.

### Theory of the firm

Ethereum is decentralized and that is fantastic. Our society has become far too dependent on untrustworthy institutions and we need to tilt the balance of power away from the few towards the many. Nevertheless, there are many, many use cases where collaboration and sharing of pooled resources lead to more resilient systems than swarms of individual laptops. A business that wants to do business on the blockchain might not want to keep its private keys on a single laptop - it probably needs something more industrial scale.

Also, in a DevOps environment, we leverage automation across the pipeline. We often have non-human actors engaged in the process of deployment and testing. Without Vault, the typical practice in the Ethereum community is to `unlock` an account for a period of time. Since there is no authentication needed to use this `unlocked` account, this creates a window of opportunity for bad actors to send transactions. A consequence of this architecture is that wallets (and private keys) are stored on personal devices - a shared wallet is pretty much impractical with conventional tools.

Having users handling passphrases with any frequency - the kind of frequency that we have in a typical development or business environment - makes exposure of passphrases likely. A tired developer will forget that they exported a variable or put a passphrase in a file.

### Vault can help

Every interaction with the Vault Ethereum backed needs to be [authenticated](https://www.vaultproject.io/docs/concepts/auth.html). Because Vault decouples authentication from storage, you can tailor the authentication mechanism to fit your needs:

* A variety of authentication providers (AppRole, AWS, Google Cloud, Kubernetes, GitHub, LDAP, MFA, Okta, RADIUS, TLS Certificates, Tokens, and Username & Password) each intended to support the unique context of the workflow.
* A sophisticated single-use token mechanism (https://www.vaultproject.io/docs/concepts/response-wrapping.html).

Every path in Vault can be protected with ACLs: You can allow some authenticated identities to import keystores, others to export them, and segregate access by account. Every access to Vault is audited as well, so it is pretty easy to diagnose access issues.

Vault encrypts all data and provides an excellent **cold storage** solution - when the Vault is sealed it requires a quorum of Shamir secret shards to bring it back online. This is functionally equivalent to what the Ethereum community would call a *multi-sig wallet*.

Furthermore, if you are an enterprise and capable of paying for [Enterprise Vault](https://www.hashicorp.com/products/vault) you can leverage HSMs as a persistence mechanism for Vault keys. This makes Vault equivalent to what the Ethereum folks call a hardware wallet. (It is very comparable to what [Gemalto and Ledger](https://www.gemalto.com/press/Pages/Gemalto-and-Ledger-Join-Forces-to-Provide--Security-Infrastructure-for-Cryptocurrency-Based-Activities-.aspx) have developed.)

## [Plugin API](https://github.com/immutability-io/vault-ethereum/blob/master/API.md)

The complete API to the plugin is documented [here](https://github.com/immutability-io/vault-ethereum/blob/master/API.md). Each API is exemplified using curl as a sample REST client.

## Plugin Setup

I assume some familiarity with Vault and Vault's plugin ecosystem. If you are not familiar, please [refer to this](https://www.vaultproject.io/guides/plugin-backends.html). I realize that it is a lot to ask for someone to be so familiar with something so new. I have a (GitHub repo that has instructions for installing Ethereum, Vault and the plugin)[https://github.com/immutability-io/immutability-project].

For this to work, you must have a Vault server already running, unsealed, and authenticated.

### Build the plugin

You can use the `Makefile` or simply us `go build` from this project's root directory.

## Install the plugin

It is assumed that your Vault configuration specifies a `plugin_directory`. Mine is:

```
$ cat vault-config.hcl

"default_lease_ttl" = "24h"

"max_lease_ttl" = "24h"

"backend" "file" {
  "path" = "/Users/immutability/etc/vault.d/data"
}

"api_addr" = "https://localhost:8200"

"listener" "tcp" {
  "address" = "localhost:8200"

  "tls_cert_file" = "/Users/immutability/etc/vault.d/vault.crt"
  "tls_client_ca_file" = "/Users/immutability/etc/vault.d/root.crt"
  "tls_key_file" = "/Users/immutability/etc/vault.d/vault.key"
}

"plugin_directory" = "/Users/immutability/etc/vault.d/vault_plugins"
```

Another configuration setting that is critical is `api_addr`. The `api_addr` must be set in order for the plugin to communicate with the Vault server during mount time.

Move the compiled plugin into Vault's configured `plugin_directory`:

```sh
$ mv vault-ethereum $HOME/etc/vault.d/vault_plugins
```

Calculate the SHA256 of the plugin and register it in Vault's plugin catalog.

```sh
$ export SHA256=$(shasum -a 256 "$HOME/etc/vault.d/vault_plugins/vault-ethereum" | cut -d' ' -f1)
$ vault write sys/plugins/catalog/ethereum-plugin \
      sha_256="${SHA256}" \
      command="vault-ethereum --ca-cert=$HOME/etc/vault.d/root.crt --client-cert=$HOME/etc/vault.d/vault.crt --client-key=$HOME/etc/vault.d/vault.key"
```

If you are using Vault in `dev` mode, you don't need to supply the certificate parameters. For any real Vault installation, however, you will be using TLS.

## Mount the Ethereum secret backend

```sh
$ vault secrets enable -path=ethereum -plugin-name=ethereum-plugin plugin
```

## Testing

I am using [Bats: Bash Automated Testing System](https://github.com/sstephenson/bats) to verify the plugin works. I recently upgraded to Vault 0.9.3 - these tests helped me discover a [problem](https://github.com/hashicorp/vault/issues/3873). I then updated the dependencies for the plugin to use [#3881](https://github.com/hashicorp/vault/pull/3881).

I have divided the tests up into 3 test cases. I initially planned to structure the tests according to the paths that are implemented by the backend:

```
Paths: framework.PathAppend(
  importPaths(&b),
  accountsPaths(&b),
  contractsPaths(&b),
),

```

However, some of the tests depend on the presence of a running Ethereum node. Also, some of these tests depend on successful mining. Therefore, I split the tests into plugin interactions that could run successfully when disconnected from an Ethereum network and tests that needed connectivity. I also created a test for plugin installation.

### Test Case: `install.bats`

With this test, we need to be a Vault administrator. Also, this test assumes that the new plugin was built and moved to the plugin directory configured (see above - `"plugin_directory" = "/Users/immutability/etc/vault.d/vault_plugins"`.)

So, assuming that you have authenticated with permissions to install the plugin, you can run the `install.bats` test case:

```
$ bats install.bats
 ✓ disable ethereum secrets plugin
 ✓ delete ethereum secrets plugin from catalog
 ✓ write ethereum secrets plugin to catalog
 ✓ enable ethereum secrets plugin

4 tests, 0 failures
```

### Test Case: `disconnected.bats`

Running the disconnected test case is simple. You need to authenticate to Vault with a policy that gives you permission to write accounts to the path where the Ethereum plugin was mounted. This policy does that. Note that this policy is very permissive. In a real use case, you would likely mount the Ethereum plugin at several paths and tightly control access within those paths:

```
path "ethereum*" {
  policy = "write"
}

```

After you have authenticated with the above permissions, you can run the `disconnected.bats` test case:

```
$ bats disconnected.bats
 ✓ list ethereum accounts - should be empty
 ✓ create test ethereum account
 ✓ read test ethereum account
 ✓ update test ethereum account no changes
 ✓ update test ethereum account blacklist
 ✓ update test ethereum account whitelist
 ✓ delete test ethereum account
 ✓ create and export test ethereum account
 ✓ import test ethereum account into test2 ethereum account
 ✓ test sign and verify
 ✓ delete test and test2 ethereum accounts

11 tests, 0 failures
```

### Test Case: `connected.bats`

To run the connected test case, you need to have access to an Ethereum network. Since this test case involves sending ETH, this better be a testnet or a private chain. This test case assumes the same private chain (`chain_id=1977`) used above. Also, as a precondition of this test case is that you have a Vault-managed Ethereum account that is funded.

If you [follow the instructions here](https://github.com/immutability-io/immutability-project), you can run this test case. Assuming that you have, you can start mining into an account using this command:

```
ETHERBASE=$(vault write -format=json ethereum/accounts/etherbase chain_id=1977 | jq .data.address | tr -d '"') ./runminer.sh etherbase
```

This will create an account named `etherbase` and pass that address to an Ethereum mining node. You need to wait some time for the account's balance to be updated.

Once you have a Vault-managed Ethereum account that is funded, you export an environment variable with this name and launch the test:

```sh
$ FUNDED_ACCOUNT=etherbase bats connected.bats
 ✓ test read etherbase balance
 ✓ test send ETH from etherbase
 ✓ test deploy contract from etherbase

3 tests, 0 failures
```

## ToDo

More (much) to come soon...

## Credits

None of this would have been possible without the fantastic [tutorial](https://www.hashicorp.com/blog/building-a-vault-secure-plugin) on Vault Plugins by Seth Vargo. Seth is one of those rare individuals who can communicate the simple essence of a complex technology in practical terms.

I had the great fortune to attend DevCon3 in November, 2017 and hear Andy Milenius speak with clarity and vision about how the Ethereum developer ecosystem should embrace the Unix philosophy - the same philosophy that makes **everything-as-code** possible: simple tools, with clear focus and purpose, driven by repeatable and interoperable mechanics. So, when I returned from DevCon3 (and dug out from my work backlog - a week away is hard) I installed `seth` and `dapp` and found inspiration.

The community chat that the [dapphub](https://dapphub.com/) guys run (esp. Andy and Mikael and Daniel Brockman) is a super warm and welcoming place that pointed me towards code that greatly helped this experiment.

## License

This code is licensed under the MPLv2 license. Please feel free to use it. Please feel free to contribute.

## Donations?

Send ETH to 0x4169c9508728285e8A9f7945D08645Bb6b3576e5 and you will be blessed in the next life.

![Donations Accepted](/doc/0x4169c9508728285e8A9f7945D08645Bb6b3576e5.png?raw=true "0x4169c9508728285e8A9f7945D08645Bb6b3576e5")
