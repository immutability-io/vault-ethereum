package ethereum

import (
	"net/rpc"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func errorPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "errors/rpc",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathErrorRPCRead,
			},
		},
		&framework.Path{
			Pattern: "errors/kill",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathErrorRPCRead,
			},
		},
	}
}

func (b *backend) pathErrorRPCRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, rpc.ErrShutdown
}
