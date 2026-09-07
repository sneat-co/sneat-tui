package commands

import (
	"errors"
	"strings"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-cli/internal/actionstore"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/spf13/cobra"
)

// validationRequiresInputMarker is the substring facade4sneatai.ErrNotCommittable
// puts in its HTTP 400 body; matching on it (rather than status code alone)
// distinguishes "you must answer a question first" from any other 400.
const validationRequiresInputMarker = "validation requires input"

// actionCommitCmd builds `sneat action commit <actionID>`. A client-held
// draft sends its current merged semantic, canonicalized (resolved values
// and display titles stripped) so the server treats it as the agent-facing
// candidate it is; a server-held action sends only the id. A refusal because
// validation still needs input exits 2 with the server's error body as
// stdout JSON, so a caller can script on the exit code alone.
func actionCommitCmd(env Env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit <actionID>",
		Short: "Commit an action",
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

			req := dto4sneatai.CommitActionRequest{
				ActionRequest: dto4sneatai.ActionRequest{
					SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID(draft.SpaceID)},
					ActionID:     actionID,
				},
			}
			if !draft.ServerHeld {
				canonical := actionspec.Canonical(draft.Semantic)
				req.Semantic = &canonical
				req.Language = draft.Language
				req.Utterance = draft.Utterance
			}
			resp, err := api.ActionCommit(cmd.Context(), req)
			if err != nil {
				if strings.Contains(err.Error(), validationRequiresInputMarker) {
					if werr := writeJSON(cmd.OutOrStdout(), map[string]string{"error": err.Error()}); werr != nil {
						return werr
					}
					return &ExitCodeError{Code: 2, Err: err}
				}
				return err
			}
			if !draft.ServerHeld {
				draft.Committed = &resp.Result
				draft.UpdatedAt = env.Now()
				if err := store.Save(draft); err != nil {
					return err
				}
			}
			return outputJSONDefault(cmd, resp, nil, nil)
		},
	}
	return cmd
}
