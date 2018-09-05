package gen

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathSignature(_ context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if err := validateFields(req, d); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	txid := d.Get("txid").(string)
	pvtkey := d.Get("pvtkey").(string)

	//var signature []string
	//signature[0]= txid
	//signature[1]= pvtkey

	signature:= []string{txid, pvtkey}

	var txSignature = strings.Join(signature, "")

	return &logical.Response{
		Data: map[string]interface{}{
			"value": txSignature,
		},
	}, nil
}