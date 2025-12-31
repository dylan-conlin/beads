package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/types"
)

func validateIssueUpdatable(id string, issue *types.Issue) error {
	if issue == nil {
		return nil
	}
	if issue.IsTemplate {
		return fmt.Errorf("Error: cannot update template %s: templates are read-only; use 'bd molecule instantiate' to create a work item", id)
	}
	return nil
}

func validateIssueClosable(id string, issue *types.Issue, force bool) error {
	if issue == nil {
		return nil
	}
	if issue.IsTemplate {
		return fmt.Errorf("Error: cannot close template %s: templates are read-only", id)
	}
	if !force && issue.Status == types.StatusPinned {
		return fmt.Errorf("Error: cannot close pinned issue %s (use --force to override)", id)
	}
	return nil
}

// validatePhaseComplete checks if an issue has a "Phase: Complete" comment.
// Returns nil if:
//   - force is true (skip validation)
//   - comments contains "Phase: Complete" (verified completion)
//   - issue is nil (for consistency with other validators)
//
// Returns error if:
//   - no "Phase: Complete" comment found (missing completion verification)
func validatePhaseComplete(id string, comments []*types.Comment, force bool) error {
	if force {
		return nil
	}
	for _, c := range comments {
		if strings.Contains(c.Text, "Phase: Complete") {
			return nil
		}
	}
	return fmt.Errorf("Error: cannot close %s: no 'Phase: Complete' comment found (use --force to override)", id)
}

func applyLabelUpdates(ctx context.Context, st storage.Storage, issueID, actor string, setLabels, addLabels, removeLabels []string) error {
	// Set labels (replaces all existing labels)
	if len(setLabels) > 0 {
		currentLabels, err := st.GetLabels(ctx, issueID)
		if err != nil {
			return err
		}
		for _, label := range currentLabels {
			if err := st.RemoveLabel(ctx, issueID, label, actor); err != nil {
				return err
			}
		}
		for _, label := range setLabels {
			if err := st.AddLabel(ctx, issueID, label, actor); err != nil {
				return err
			}
		}
	}

	// Add labels
	for _, label := range addLabels {
		if err := st.AddLabel(ctx, issueID, label, actor); err != nil {
			return err
		}
	}

	// Remove labels
	for _, label := range removeLabels {
		if err := st.RemoveLabel(ctx, issueID, label, actor); err != nil {
			return err
		}
	}

	return nil
}
