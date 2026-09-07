package commands

import (
	"encoding/json"
	"fmt"

	"github.com/sneat-co/sneat-ai-backend/dto4sneatai"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/spf13/cobra"
)

// Query builds the top-level `sneat query` command: "what to buy for X /
// before Y", answered from the Space's lists and happenings.
func Query(env Env) *cobra.Command {
	var space, jsonBody string
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Ask what to buy/schedule for contacts before a horizon",
		RunE: func(cmd *cobra.Command, _ []string) error {
			spaceID, err := resolveSpaceID(cmd, env, space)
			if err != nil {
				return err
			}
			if jsonBody == "" {
				return fmt.Errorf(`--json is required, e.g. '{"contacts":[{"mention":"Vasilisa"}],"before":{"relative":"end_of_month"}}'`)
			}
			var req dto4sneatai.QueryRequest
			if err := json.Unmarshal([]byte(jsonBody), &req); err != nil {
				return fmt.Errorf("--json: %w", err)
			}
			req.SpaceRequest = dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID(spaceID)}
			api, err := env.NewActionsAPI(configFromCmd(cmd, env.Getenv))
			if err != nil {
				return err
			}
			resp, err := api.Query(cmd.Context(), req)
			if err != nil {
				return err
			}
			return outputJSONDefault(cmd, resp, nil, nil)
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "space id, or 'family'/'private' (default: current space or family)")
	cmd.Flags().StringVar(&jsonBody, "json", "", `query body as JSON: {"contacts":[...],"before":{...},"listID":"..."}`)
	return cmd
}
