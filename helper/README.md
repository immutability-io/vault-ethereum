# Helpers to get up and running

In this directory are a few scripts to help you get up and running with vault and the ethereum plugin. These will install vault and the plugin, initialize vault and configure the plugin. Vault tokens and key shards will be encrypted on the filesystem.

## Prerequisites

You need to [install Keybase](https://keybase.io/) and create an identity in Keybase. Your Keybase identity will be used to encrypt the root token and keyshards. 

I am also using the fantastic [jq](https://stedolan.github.io/jq/) because you don't ever want to use JSON over REST without it. All vault commands support `-format=json` and this allows you to pipe output directly into `jq` to access specific keys.

## Install vault and plugin

Installing vault so that it uses TLS is a little more complicated than merely downloading vault. You need to generate a CA cert, sign a CSR and generate a vault configuration. Ideally, you would check the signature of the vault executable. Installing the plugin involves creating a plugins directory, configuring vault to know where that directory is, and similarly checking the signature of the plugin. This script does all of this for you. 

### install_vault.sh

The usage can be seen by executing the script without any parameters:

```
$ ./install_vault.sh
Usage: bash install_vault.sh OPTIONS

OPTIONS:
  --linux	Install Linux version
  --darwin	Install Darwin (MacOS) version
```

To install on MacOS: `./install_vault.sh --darwin`. This will create a directory: `$HOME/etc/vault.d`. It will put the vault executable in `/usr/local/bin/`.

## Initialize vault

Vault needs to be initialized before use. This involves generating the key shards that are used to unseal vault and generate the master encryption key. Our script also unseals vault. All key shards and the root token are encrypted using the supplied Keybase identity.

### initialize_vault.sh

The usage can be seen by executing the script without any parameters:

```
$ ./initialize_vault.sh
Usage: bash initialize_vault.sh OPTIONS

OPTIONS:
  [keybase]	Name of Keybase user to encrypt Vault keys with
```

To initialize vault using the Keybase identity `cypherhat`:

```
$ ./initialize_vault.sh cypherhat
Key                Value
---                -----
Seal Type          shamir
Sealed             true
Total Shares       5
Threshold          3
Unseal Progress    1/3
Unseal Nonce       6bf3d366-0fb9-dce8-bd98-5a2e8a2bff8b
Version            0.11.0
HA Enabled         false
Key                Value
---                -----
Seal Type          shamir
Sealed             true
Total Shares       5
Threshold          3
Unseal Progress    2/3
Unseal Nonce       6bf3d366-0fb9-dce8-bd98-5a2e8a2bff8b
Version            0.11.0
HA Enabled         false
Key             Value
---             -----
Seal Type       shamir
Sealed          false
Total Shares    5
Threshold       3
Version         0.11.0
Cluster Name    vault-cluster-4ccd2d1e
Cluster ID      9ed36d67-cdb2-d65c-b0c8-7dea61848668
HA Enabled      false
Key             Value
---             -----
Seal Type       shamir
Sealed          false
Total Shares    5
Threshold       3
Version         0.11.0
Cluster Name    vault-cluster-4ccd2d1e
Cluster ID      9ed36d67-cdb2-d65c-b0c8-7dea61848668
HA Enabled      false
Key             Value
---             -----
Seal Type       shamir
Sealed          false
Total Shares    5
Threshold       3
Version         0.11.0
Cluster Name    vault-cluster-4ccd2d1e
Cluster ID      9ed36d67-cdb2-d65c-b0c8-7dea61848668
HA Enabled      false
```

If we look at the resulting file system, we see that all secrets (encrypted) are named using the Keybase identity:

```
$ ls -ltr cypherhat_*
-rw-r--r--  1 cypherhat  staff  1742 Sep  1 08:34 cypherhat_VAULT_ROOT_TOKEN.txt
-rw-r--r--  1 cypherhat  staff  1786 Sep  1 08:34 cypherhat_UNSEAL_0.txt
-rw-r--r--  1 cypherhat  staff  1786 Sep  1 08:34 cypherhat_UNSEAL_1.txt
-rw-r--r--  1 cypherhat  staff  1786 Sep  1 08:35 cypherhat_UNSEAL_2.txt
-rw-r--r--  1 cypherhat  staff  1786 Sep  1 08:35 cypherhat_UNSEAL_3.txt
-rw-r--r--  1 cypherhat  staff  1786 Sep  1 08:35 cypherhat_UNSEAL_4.txt
```

## Configure ethereum plugin

Before we configure the plugin, we have to log into vault as an administrator. We will use the root token to do this. Since the root token is encrypted using Keybase, we have to login to Keybase first:

![Keybase after login as cypherhat](../docs/keybase.png?raw=true "Keybase Logged In")

The usage can be seen by executing the script without any parameters:

```
$ ./config_plugin.sh
Usage: bash config_plugin.sh OPTIONS

OPTIONS:
  [keybase]	Name of Keybase user used to encrypt Vault keys
```

The `config_plugin.sh` script will authenticate to Vault using the following approach:

```
$ export VAULT_TOKEN=$(keybase decrypt -i cypherhat_VAULT_ROOT_TOKEN.txt)
Message authored by cypherhat
```

Assuming this works (which it should if you are logged into Keybase as the user - in my case that user is `cypherhat` - you specify to the script), you should see something like:

```
$ ./config_plugin.sh cypherhat
Message authored by cypherhat
ADDING TO CATALOG: sys/plugins/catalog/ethereum-plugin
Success! Data written to: sys/plugins/catalog/ethereum-plugin
MOUNTING: ethereum/mainnet
Success! Enabled the ethereum-plugin plugin at: ethereum/mainnet/
MOUNTING: ethereum/rinkeby
Success! Enabled the ethereum-plugin plugin at: ethereum/rinkeby/
CONFIGURE: ethereum/mainnet
Key                Value
---                -----
api_key            n/a
bound_cidr_list    <nil>
chain_id           1
rpc_url            https://mainnet.infura.io
CONFIGURE: ethereum/rinkeby
Key                Value
---                -----
api_key            n/a
bound_cidr_list    <nil>
chain_id           4
rpc_url            https://rinkeby.infura.io
```

This will mount the plugin at 2 paths:

- [x] - `ethereum/mainnet/` which points to the live mainnet (using the Infura endpoint)
- [x] - `ethereum/rinkeby/` which points to the Rinkeby testnet (using the Infura endpoint)

## Using the plugin

If you scrape your environment for variables with `VAULT`, you should see something like:

```
$ env | grep VAULT
VAULT_ADDR=https://localhost:8200
VAULT_CACERT=/Users/cypherhat/etc/vault.d/root.crt
```

You can play with the plugin before you authenticate using the unauthenticated paths. For example, you can convert units of ethereum:

```
$ vault write -format=json ethereum/mainnet/convert amount=12 unit_to=babbage unit_from=wei | jq .data
{
  "amount_from": "12",
  "amount_to": "0.000012",
  "unit_from": "wei",
  "unit_to": "mwei"
}
```


However if you want to access features that will use private keys, you still have to authenticate to Vault. 

```
$ vault write -f -format=json ethereum/rinkeby/accounts/muchwow | jq .data
Error writing data to ethereum/rinkeby/accounts/muchwow: Error making API request.

URL: PUT https://localhost:8200/v1/ethereum/rinkeby/accounts/muchwow
Code: 400. Errors:

* missing client token
```

To use the plugin in anything like a production setting, you will want to create policies and attach them to various identities to allow the kind of access you wish. However, for this exercise, you can authenticate with the root token:

```
$ export VAULT_TOKEN=$(keybase decrypt -i cypherhat_VAULT_ROOT_TOKEN.txt)
Message authored by cypherhat
```

Now you can create accounts:

```
$ vault write -f -format=json ethereum/rinkeby/accounts/muchwow | jq .data
{
  "address": "0xd90b08955547e97e325a17ae223f3f482e9e0a37",
  "blacklist": null,
  "spending_limit_total": "0",
  "spending_limit_tx": "0",
  "total_spend": "0",
  "whitelist": null
}
```

## Don't delete your encrypted keys!!!

You can move your Vault data and keys to *cold* storage by terminating the vault server and moving the data and encrypted keys to offline storage - like an optical disk or a flash drive.

I've mounted a flash drive: `/Volumes/cold`. Running the script `cold.sh` will terminate the vault server and move the data and encrypted keys to offline storage:

```
$ ./cold.sh
Usage: bash cold.sh OPTIONS

OPTIONS:
  [keybase]	Name of Keybase user used to encrypt Vault keys
  [path]	Path to mounted Flash drive or other media

$ ./cold.sh cypherhat /Volumes/cold
```

After you run this, there is nothing left on the original file system containing your private keys. **You should always logout of Keybase after doing this - the script does not do this.**

To restore from cold storage, mount the flash drive, login to Keybase and run the `hot.sh` script. This will start the vault server and unseal it:

```
$ ./hot.sh
Usage: bash hot.sh OPTIONS

OPTIONS:
  [keybase]	Name of Keybase user used to encrypt Vault keys
  [path]	Path to cold storage

  
$ ./hot.sh cypherhat /Volumes/cold
Message authored by cypherhat
Key                Value
---                -----
Seal Type          shamir
Sealed             true
Total Shares       5
Threshold          3
Unseal Progress    1/3
Unseal Nonce       04c79e48-3c16-1d73-0a6b-1c1ce16a8052
Version            0.11.0
HA Enabled         false
Message authored by cypherhat
Key                Value
---                -----
Seal Type          shamir
Sealed             true
Total Shares       5
Threshold          3
Unseal Progress    2/3
Unseal Nonce       04c79e48-3c16-1d73-0a6b-1c1ce16a8052
Version            0.11.0
HA Enabled         false
Message authored by cypherhat
Key             Value
---             -----
Seal Type       shamir
Sealed          false
Total Shares    5
Threshold       3
Version         0.11.0
Cluster Name    vault-cluster-4ccd2d1e
Cluster ID      9ed36d67-cdb2-d65c-b0c8-7dea61848668
HA Enabled      false
```

All that remains is for you to authenticate to vault to perform whatever actions you wish.