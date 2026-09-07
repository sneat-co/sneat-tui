package commands

import (
	"errors"

	"github.com/sneat-co/sneat-cli/internal/actionstore"
	"github.com/spf13/cobra"
)

// actionShowCmd builds `sneat action show <actionID>`: the local draft for a
// client-held action, or action_get for a server-held one.
func actionShowCmd(env Env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <actionID>",
		Short: "Show an action (local draft, or fetched from the server)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			actionID := args[0]
			store, err := newActionStore(env)
			if err != nil {
				return err
			}
			draft, err := store.Load(actionID)
			if err != nil {
				if errors.Is(err, actionstore.ErrNotFound) {
					return errNoSuchDraft(actionID)
				}
				return err
			}
			if !draft.ServerHeld {
				return outputJSONDefault(cmd, draft, nil, nil)
			}
			api, err := env.NewActionsAPI(configFromCmd(cmd, env.Getenv))
			if err != nil {
				return err
			}
			resp, err := api.ActionGet(cmd.Context(), draft.SpaceID, actionID)
			if err != nil {
				return err
			}
			return outputJSONDefault(cmd, resp, nil, nil)
		},
	}
	return cmd
}
