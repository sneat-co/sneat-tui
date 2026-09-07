package commands

import (
	"github.com/sneat-co/sneat-cli/internal/actionstore"
	"github.com/spf13/cobra"
)

var actionListHeaders = []string{"ID", "SPACE", "STATUS", "KIND", "UPDATED"}

// actionListCmd builds `sneat action list`: a summary of local drafts (both
// client-held drafts and locally-remembered server-held index entries).
func actionListCmd(env Env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List local action drafts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			store, err := newActionStore(env)
			if err != nil {
				return err
			}
			drafts, err := store.List()
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(drafts))
			for _, d := range drafts {
				rows = append(rows, actionListRow(d))
			}
			return outputJSONDefault(cmd, drafts, actionListHeaders, rows)
		},
	}
	return cmd
}

func actionListRow(d actionstore.Draft) []string {
	status := "draft"
	switch {
	case d.Cancelled:
		status = "cancelled"
	case d.Committed != nil:
		status = "committed"
	case d.ServerHeld:
		status = "server-held"
	case d.Validation.CanCommit:
		status = "validated"
	}
	return []string{d.ActionID, d.SpaceID, status, string(d.Semantic.Kind), d.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")}
}
