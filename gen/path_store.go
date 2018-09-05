package gen

import (
	"context"
	"net/http"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type Key struct {
	PrivateKey string `json:"private_key"`
}

func (b *backend) storeKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := validateFields(req, d); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	pvtkey := d.Get("pvtkey").(string)

	pvtkeyJSON := Key{
		PrivateKey: pvtkey,
	}

	entry, err := logical.StorageEntryJSON(req.Path, pvtkeyJSON)

	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)

	return &logical.Response{
		Data: map[string]interface{}{
			"pvtKey": pvtkeyJSON.PrivateKey,
		},
	}, nil
}