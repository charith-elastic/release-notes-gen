// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package main

import (
	"fmt"
	"os"

	"github.com/elastic/release-notes-gen/config"
	"github.com/elastic/release-notes-gen/github"
	"github.com/spf13/cobra"
)

var (
	appCommit  = "unknown"
	appVersion = "0.1.0"

	args = config.Args{}
)

func main() {
	cmd := &cobra.Command{
		Use:           "release-notes-gen",
		Short:         "Generate release notes from GitHub pull requests",
		Example:       "GITHUB_TOKEN=xxxyyy ./release-notes-gen --conf=example/config.yaml --template=example/template.tpl --label=v1.2.0",
		Version:       fmt.Sprintf("%s (%s)", appVersion, appCommit),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE:          doRun,
	}

	cmd.Flags().StringVar(&args.ConfPath, "conf", "", "Path to the configuration file")
	_ = cmd.MarkFlagRequired("conf")

	cmd.Flags().StringVar(&args.FilterLabel, "label", "", "Label to filter PRs by (e.g. v1.2.0)")
	_ = cmd.MarkFlagRequired("label")

	cmd.Flags().StringVar(&args.TemplatePath, "template", "", "Path to the release notes template")
	_ = cmd.MarkFlagRequired("template")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func doRun(_ *cobra.Command, _ []string) error {
	conf, err := config.Load(args, os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return err
	}

	prs, err := github.LoadPullRequests(conf)
	if err != nil {
		return err
	}

	if len(prs) == 0 {
		return fmt.Errorf("no pull requests found matching label %s", conf.FilterLabel)
	}

	groupedPRs := groupPullRequests(conf, prs)
	if err := render(conf, groupedPRs); err != nil {
		return err
	}

	return nil
}

func groupPullRequests(conf *config.Config, prs []github.PullRequest) map[string][]github.PullRequest {
	groups := make(map[string][]github.PullRequest)

PR_LOOP:
	for _, pr := range prs {
		for _, lbl := range conf.GroupOrder {
			if _, ok := pr.Labels[lbl]; ok {
				groups[lbl] = append(groups[lbl], pr)
				continue PR_LOOP
			}
		}

		groups[config.NoGroupKey] = append(groups[config.NoGroupKey], pr)
	}

	return groups
}

func render(conf *config.Config, groups map[string][]github.PullRequest) error {
	params := struct {
		FilterLabel string
		Repo        string
		Groups      map[string][]github.PullRequest
		GroupLabels map[string]string
		GroupOrder  []string
	}{
		FilterLabel: conf.FilterLabel,
		Repo:        conf.Repository,
		Groups:      groups,
		GroupLabels: conf.GroupLabels,
		GroupOrder:  conf.GroupOrder,
	}

	return conf.Template.Execute(os.Stdout, params)
}
