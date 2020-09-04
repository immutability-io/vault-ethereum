# Vault Ethereum Plugin v0.3.0

The first incarnation of the `vault-ethereum` plugin was an exercise in [experimenting with an idea](https://www.hashicorp.com/resources/vault-platform-enterprise-blockchain) and [proving a point](https://immutability.io/). 2 years later, I feel both ends were acheived.

Having had several occasions to take this PoC to production with companies in the financial and blockchain communities [(plug for Immutability, LLC's custom development!)](mailto:jeff@immutability.io) I've decided to release an upgrade that tries to make the development experience better. I've also restricted the surface area of the plugin to a minimum.

## Testing - in one terminal...

```sh

$ cd $GOPATH/src/github.com/immutability-io/vault-ethereum
$ make docker-build
$ make run

```

Then, **open a different terminal**...

```sh

$ cd $GOPATH/src/github.com/immutability-io/vault-ethereum/docker

# Authenticate
$ source ./local-test.sh auth
$ ./demo.sh > README.md

```

## View the demo

If everything worked... And you have run the command above, your demo is had by viewing the results: 

```sh
$ cat ./README.md
```

If everything didn't work, tell me why.

## What is the API?

The best way to understand the API is to use the `path-help` command. For example:

```sh
$ vault path-help vault-ethereum/accounts/bob/deploy                                                                [±new-version ●]
Request:        accounts/bob/deploy
Matching Route: ^accounts/(?P<name>\w(([\w-.]+)?\w)?)/deploy$

Deploy a smart contract from an account.

## PARAMETERS

    abi (string)

        The contract ABI.

    address (string)

        <no description>

    bin (string)

        The compiled smart contract.

    gas_limit (string)

        The gas limit for the transaction - defaults to 0 meaning estimate.

    name (string)

        <no description>

    version (string)

        The smart contract version.

## DESCRIPTION

Deploy a smart contract to the network.
```

## I still need help

[Please reach out to me](mailto:jeff@immutability.io). 
