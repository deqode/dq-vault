package api

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	_, err := req.Storage.List(ctx, "users/")

	return nil, err

}
