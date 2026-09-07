package commands

import (
	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/spf13/cobra"
)

// actionSchemaCmd builds `sneat action schema`: the machine-readable
// description of the semantic model, fetched from the server by default (so
// a running agent always sees what this deployment actually accepts) or
// printed from the compiled-in actionspec module with --offline.
func actionSchemaCmd(env Env) *cobra.Command {
	var offline bool
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Print the Action Protocol's machine-readable semantic schema",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if offline {
				return outputJSONDefault(cmd, actionspec.CurrentSchema(), nil, nil)
			}
			api, err := env.NewActionsAPI(configFromCmd(cmd, env.Getenv))
			if err != nil {
				return err
			}
			schema, err := api.Schema(cmd.Context())
			if err != nil {
				return err
			}
			return outputJSONDefault(cmd, schema, nil, nil)
		},
	}
	cmd.Flags().BoolVar(&offline, "offline", false, "print the schema compiled into this CLI build instead of fetching it")
	return cmd
}
