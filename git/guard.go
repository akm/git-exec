package git

import (
	"fmt"
	"strings"
)

type GuardResult struct {
	UncommittedChanges string
	UntrackedFiles     string
	Skipped            bool
}

func (g *GuardResult) Message() string {
	var r string
	if len(g.UncommittedChanges) > 0 && len(g.UntrackedFiles) > 0 {
		r = "There are uncommitted changes and untracked files"
	} else if len(g.UntrackedFiles) > 0 {
		r = "There are untracked files"
	} else {
		r = "There are uncommitted changes"
	}
	if g.Skipped {
		r += " but guard was skipped by options"
	}
	return r
}

func (g *GuardResult) Format() string {
	parts := []string{g.Message()}
	if len(g.UncommittedChanges) > 0 {
		parts = append(parts, fmt.Sprintf("Uncommitted changes:\n%s", g.UncommittedChanges))
	}
	if len(g.UntrackedFiles) > 0 {
		parts = append(parts, fmt.Sprintf("Untracked files:\n%s", g.UntrackedFiles))
	}
	return strings.Join(parts, "\n\n")
}

type GuardOptions struct {
	SkipGuard                   bool
	SkipGuardUncommittedChanges bool
	SkipGuardUntrackedFiles     bool
}

func Guard(opts *GuardOptions) (*GuardResult, error) {
	diff, err := UncommittedChanges()
	if err != nil {
		return nil, err
	}

	untrackedFiles, err := UntrackedFiles()
	if err != nil {
		return nil, err
	}

	if len(diff) == 0 && len(untrackedFiles) == 0 {
		return nil, nil
	}

	return &GuardResult{
		UncommittedChanges: diff,
		UntrackedFiles:     untrackedFiles,
		Skipped: opts.SkipGuard ||
			(opts.SkipGuardUncommittedChanges && len(diff) > 0) ||
			(opts.SkipGuardUntrackedFiles && len(untrackedFiles) > 0),
	}, nil
}
