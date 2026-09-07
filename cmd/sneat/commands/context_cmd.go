package commands

import (
	"github.com/spf13/cobra"
)

// Context builds the top-level `sneat context` command: the Space snapshot
// (contacts, lists, happenings, now/today/weekday) an agent reads before
// deciding what an utterance means.
func Context(env Env) *cobra.Command {
	var space string
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Print a Space snapshot for an agent to read before acting",
		RunE: func(cmd *cobra.Command, _ []string) error {
			spaceID, err := resolveSpaceID(cmd, env, space)
			if err != nil {
				return err
			}
			api, err := env.NewActionsAPI(configFromCmd(cmd, env.Getenv))
			if err != nil {
				return err
			}
			resp, err := api.Context(cmd.Context(), spaceID)
			if err != nil {
				return err
			}
			return outputJSONDefault(cmd, resp, nil, nil)
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "space id, or 'family'/'private' (default: current space or family)")
	return cmd
}
