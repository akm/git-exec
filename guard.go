package main

import "fmt"

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
	return fmt.Sprintf("%s\nUncommitted changes:\n%s\n\nUntracked files:\n%s\n",
		g.Message(),
		g.uncommittedChanges,
		g.untrackedFiles,
	)
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
