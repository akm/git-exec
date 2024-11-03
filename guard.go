package main

import (
	"fmt"
	"strings"
)

type guardResult struct {
	uncommittedChanges string
	untrackedFiles     string
	skipped            bool
}

func (g *guardResult) Message() string {
	var r string
	if len(g.uncommittedChanges) > 0 && len(g.untrackedFiles) > 0 {
		r = "There are uncommitted changes and untracked files"
	} else if len(g.untrackedFiles) > 0 {
		r = "There are uncommitted changes"
	} else {
		r = "There are untracked files"
	}
	if g.skipped {
		r += " but guard was skipped by options"
	}
	return r
}

func (g *guardResult) Format() string {
	parts := []string{g.Message()}
	if len(g.uncommittedChanges) > 0 {
		parts = append(parts, fmt.Sprintf("Uncommitted changes:\n%s", g.uncommittedChanges))
	}
	if len(g.untrackedFiles) > 0 {
		parts = append(parts, fmt.Sprintf("Untracked files:\n%s", g.untrackedFiles))
	}
	return strings.Join(parts, "\n\n")
}

func guard(opts *Options) (*guardResult, error) {
	diff, err := uncommittedChanges()
	if err != nil {
		return nil, err
	}

	untrackedFiles, err := untrackedFiles()
	if err != nil {
		return nil, err
	}

	if len(diff) == 0 && len(untrackedFiles) == 0 {
		return nil, nil
	}

	return &guardResult{
		uncommittedChanges: diff,
		untrackedFiles:     untrackedFiles,
		skipped: opts.SkipGuard ||
			(opts.SkipGuardUncommittedChanges && len(diff) > 0) ||
			(opts.SkipGuardUntrackedFiles && len(untrackedFiles) > 0),
	}, nil
}
