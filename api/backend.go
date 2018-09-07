package api

import (
	"context"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/pkg/errors"
)

// Factory creates a new usable instance of this secrets engine.
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(c)
	if err := b.Setup(ctx, c); err != nil {
		return nil, errors.Wrap(err, "failed to create factory")
	}
	return b, nil
}

// backend is the actual backend.
type backend struct {
	*framework.Backend
	logger log.Logger
}

// Backend creates a new backend.
func Backend(c *logical.BackendConfig) *backend {
	var b backend

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help:        backendHelp,
		Paths: []*framework.Path{

			// api/register
			&framework.Path{
				Pattern:      "register",
				HelpSynopsis: "Registers a new user in vault with mnemonic and UUID",
				HelpDescription: `

Registers new user in vault using UUID. Generates mnemonics if not provided and store it in vault.

`,
				Fields: map[string]*framework.FieldSchema{
					"uuid": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "UUID of new user",
						Default:     "",
					},
					"mnemonic": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Mnemonic for bip39 seed",
						Default:     "",
					},
					"passphrase": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Passphrase for bip39 seed",
						Default:     "",
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.pathRegister,
				},
			},

			// api/signature
			&framework.Path{
				Pattern:         "signature",
				HelpSynopsis:    "Generate a signature",
				HelpDescription: "Generates a signature from stored mnemonic and passphrase using deviation path",
				Fields: map[string]*framework.FieldSchema{
					"uuid": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "UUID of user to read credentials",
					},
					"path": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Deviation path to obtain keys",
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.pathSignature,
				},
			},

			// api/info
			&framework.Path{
				Pattern:      "info",
				HelpSynopsis: "Display information about this plugin",
				HelpDescription: `

Displays information about the plugin, such as the plugin version and where to
get help.

`,
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.pathInfo,
				},
			},
		},
	}
	return &b
}

const backendHelp = `
The gen secrets engine generates passwords and passphrases, and optionally
stores the resulting password in an accessor.
`
