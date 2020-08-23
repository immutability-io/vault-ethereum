# Vault Ethereum Plugin v0.3.0

The first incarnation of the `vault-ethereum` plugin was an exercise in [exposing an idea](https://www.hashicorp.com/resources/vault-platform-enterprise-blockchain) and [proving a point](https://immutability.io/). 2 years later, I feel both ends were acheived.

Having had several occasions to take this PoC to production with companies in the financial community [(plug for Immutability, LLC's custom development!)](mailto:jeff@immutability.io) I've decided to release an upgrade that tries to make the development experience better. I've also restricted the surface area of the plugin to a minimum, **in hope of driving a standard wallet interact for enterprises** as well as to reduce the threat surface.

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

If everything worked... And you have run the command above, [this is your demo](./docker/README.md).

If everything didn't work, tell me why.