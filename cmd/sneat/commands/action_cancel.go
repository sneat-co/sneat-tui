package commands

import (
	"errors"

	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-cli/internal/actionstore"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/spf13/cobra"
)

// actionCancelCmd builds `sneat action cancel <actionID>`. A client-held
// draft is marked cancelled and its file deleted; a server-held action is
// cancelled via action_cancel.
func actionCancelCmd(env Env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel <actionID>",
		Short: "Cancel an action",
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

			if draft.ServerHeld {
				api, err := env.NewActionsAPI(configFromCmd(cmd, env.Getenv))
				if err != nil {
					return err
				}
				if err := api.ActionCancel(cmd.Context(), dto4sneatai.ActionRequest{
					SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID(draft.SpaceID)},
					ActionID:     actionID,
				}); err != nil {
					return err
				}
				if err := store.Delete(actionID); err != nil {
					return err
				}
				return outputJSONDefault(cmd, map[string]string{"actionID": actionID, "status": "cancelled"}, nil, nil)
			}

			draft.Cancelled = true
			draft.UpdatedAt = env.Now()
			if err := store.Delete(actionID); err != nil {
				return err
			}
			return outputJSONDefault(cmd, draft, nil, nil)
		},
	}
	return cmd
}
