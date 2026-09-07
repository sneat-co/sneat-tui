package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-cli/internal/actionstore"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/spf13/cobra"
)

// actionAddCmd builds `sneat action add <actionID>`, applying a semantic
// patch to a pending action — merged locally and re-validated for a
// client-held draft, or sent to action_patch for a server-held one.
func actionAddCmd(env Env) *cobra.Command {
	var jsonPatch, utterance, language string
	cmd := &cobra.Command{
		Use:   "add <actionID>",
		Short: "Apply a semantic patch to a pending action",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			actionID := args[0]
			if jsonPatch == "" {
				return fmt.Errorf("--json is required: a semantic patch, e.g. '{\"schedule\":{...}}'")
			}
			var patch actionspec.Semantic
			if err := json.Unmarshal([]byte(jsonPatch), &patch); err != nil {
				return fmt.Errorf("--json: %w", err)
			}
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

			if draft.ServerHeld {
				resp, err := api.ActionPatch(cmd.Context(), dto4sneatai.PatchActionRequest{
					ActionRequest: dto4sneatai.ActionRequest{
						SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID(draft.SpaceID)},
						ActionID:     actionID,
					},
					Utterance: utterance,
					Language:  language,
					Patch:     patch,
				})
				if err != nil {
					return err
				}
				return outputJSONDefault(cmd, resp, nil, nil)
			}

			// The resolved semantic from the last validate is the merge base so
			// already-resolved contact IDs and dates are kept rather than
			// re-guessed from raw mentions.
			merged := actionspec.Merge(draft.Semantic, patch)
			var committedHappeningID string
			if draft.Committed != nil {
				committedHappeningID = draft.Committed.HappeningID
			}
			resp, err := api.ActionValidate(cmd.Context(), dto4sneatai.ValidateActionRequest{
				SpaceRequest:         dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID(draft.SpaceID)},
				ActionID:             actionID,
				Semantic:             &merged,
				Language:             language,
				Utterance:            utterance,
				CommittedHappeningID: committedHappeningID,
			})
			if err != nil {
				return err
			}
			if language != "" {
				draft.Language = language
			}
			if utterance != "" {
				draft.Utterance = utterance
			}
			draft.Semantic = resp.Semantic
			draft.Validation = resp.Validation
			draft.UpdatedAt = env.Now()
			if err := store.Save(draft); err != nil {
				return err
			}
			resp.ActionID = actionID
			return outputJSONDefault(cmd, resp, nil, nil)
		},
	}
	f := cmd.Flags()
	f.StringVar(&jsonPatch, "json", "", "semantic patch as JSON (required)")
	f.StringVar(&utterance, "utterance", "", "the user's follow-up utterance")
	f.StringVar(&language, "language", "", "utterance language: en|ru")
	return cmd
}
