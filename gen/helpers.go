package gen

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// Node -- stores data related of a node derived from master node (HD Wallet)
type node struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Path       string `json:"path"`
	Address    string `json:"address"`
}

// CheckError checks for any potential errors
func checkError(err error, message string) {
	if err != nil {
		log.Fatalf("%v - %v", message, err)
	}
}

// errMissingField returns a logical response error that prints a consistent
// error message for when a required field is missing.
func errMissingField(field string) *logical.Response {
	return logical.ErrorResponse(fmt.Sprintf("missing required field '%s'", field))
}

// validationErr returns an error that corresponds to a validation error.
func validationErr(msg string) error {
	return logical.CodedError(http.StatusUnprocessableEntity, msg)
}

// validateFields verifies that no bad arguments were given to the request.
func validateFields(req *logical.Request, data *framework.FieldData) error {
	var unknownFields []string
	for k := range req.Data {
		if _, ok := data.Schema[k]; !ok {
			unknownFields = append(unknownFields, k)
		}
	}

	if len(unknownFields) > 0 {
		// Sort since this is a human error
		sort.Strings(unknownFields)

		return fmt.Errorf("unknown fields: %q", unknownFields)
	}

	return nil
}
