package commands

import (
	"encoding/json"
	"fmt"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-cli/internal/actionstore"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/spf13/cobra"
)

// actionNewCmd builds `sneat action new`, starting an action from a semantic
// candidate. By default it holds the draft client-side (a local file plus a
// stateless action_validate call); --server instead persists it on the
// server via action_create.
func actionNewCmd(env Env) *cobra.Command {
	var space, jsonSemantic, utterance, language string
	var server bool
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Start an action from a semantic candidate (JSON)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			spaceID, err := resolveSpaceID(cmd, env, space)
			if err != nil {
				return err
			}
			var semantic actionspec.Semantic
			if jsonSemantic == "" {
				return fmt.Errorf("--json is required: a semantic candidate, e.g. '{\"kind\":\"buy\",...}'")
			}
			if err := json.Unmarshal([]byte(jsonSemantic), &semantic); err != nil {
				return fmt.Errorf("--json: %w", err)
			}
			api, err := env.NewActionsAPI(configFromCmd(cmd, env.Getenv))
			if err != nil {
				return err
			}
			now := env.Now()

			if server {
				resp, err := api.ActionCreate(cmd.Context(), dto4sneatai.CreateActionRequest{
					SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID(spaceID)},
					Utterance:    utterance,
					Language:     language,
					Semantic:     semantic,
				})
				if err != nil {
					return err
				}
				store, err := newActionStore(env)
				if err != nil {
					return err
				}
				if err := store.Save(actionstore.Draft{
					ActionID: resp.ActionID, SpaceID: spaceID, ServerHeld: true,
					CreatedAt: now, UpdatedAt: now,
				}); err != nil {
					return err
				}
				return outputJSONDefault(cmd, resp, nil, nil)
			}

			actionID, err := newActionID()
			if err != nil {
				return err
			}
			resp, err := api.ActionValidate(cmd.Context(), dto4sneatai.ValidateActionRequest{
				SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID(spaceID)},
				ActionID:     actionID,
				Semantic:     &semantic,
				Language:     language,
				Utterance:    utterance,
			})
			if err != nil {
				return err
			}
			store, err := newActionStore(env)
			if err != nil {
				return err
			}
			draft := actionstore.Draft{
				ActionID: actionID, SpaceID: spaceID, Language: language, Utterance: utterance,
				Semantic: resp.Semantic, Validation: resp.Validation,
				CreatedAt: now, UpdatedAt: now,
			}
			if err := store.Save(draft); err != nil {
				return err
			}
			resp.ActionID = actionID
			return outputJSONDefault(cmd, resp, nil, nil)
		},
	}
	f := cmd.Flags()
	f.StringVar(&space, "space", "", "space id, or 'family'/'private' (default: current space or family)")
	f.StringVar(&jsonSemantic, "json", "", "semantic candidate as JSON (required)")
	f.StringVar(&utterance, "utterance", "", "the user's original utterance (audit + language detection)")
	f.StringVar(&language, "language", "", "utterance language: en|ru")
	f.BoolVar(&server, "server", false, "hold the action on the server (action_create) instead of as a local draft")
	return cmd
}
