package commands

import (
	"errors"

	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-cli/internal/actionstore"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/spf13/cobra"
)

// actionValidateCmd builds `sneat action validate <actionID>`, re-running
// validation without changing the semantic candidate: inline for a
// client-held draft, by id for a server-held action.
func actionValidateCmd(env Env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate <actionID>",
		Short: "Re-validate an action",
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
			api, err := env.NewActionsAPI(configFromCmd(cmd, env.Getenv))
			if err != nil {
				return err
			}

			req := dto4sneatai.ValidateActionRequest{
				SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID(draft.SpaceID)},
				ActionID:     actionID,
			}
			if !draft.ServerHeld {
				semantic := draft.Semantic
				req.Semantic = &semantic
				req.Language = draft.Language
				if draft.Committed != nil {
					req.CommittedHappeningID = draft.Committed.HappeningID
				}
			}
			resp, err := api.ActionValidate(cmd.Context(), req)
			if err != nil {
				return err
			}
			if !draft.ServerHeld {
				draft.Semantic = resp.Semantic
				draft.Validation = resp.Validation
				draft.UpdatedAt = env.Now()
				if err := store.Save(draft); err != nil {
					return err
				}
				resp.ActionID = actionID
			}
			return outputJSONDefault(cmd, resp, nil, nil)
		},
	}
	return cmd
}
