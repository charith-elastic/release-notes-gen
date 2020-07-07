// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
)

var errNoToken = errors.New("token not found: provide a valid GitHub token using the GITHUB_TOKEN environment variable")

const (
	NoGroupKey     = "nogroup"
	noGroupDisplay = "No group"
)

// Args holds command-line arguments passed to the application.
type Args struct {
	ConfPath     string
	FilterLabel  string
	TemplatePath string
}

// Config holds configuration parameters required by the application.
type Config struct {
	FilterLabel  string
	GitHubToken  string
	GroupLabels  map[string]string
	GroupOrder   []string
	IgnoreLabels map[string]struct{}
	Repository   string
	Template     *template.Template
}

type configFile struct {
	Repository string `json:"repository"`
	Groups     []struct {
		Label   string `json:"label"`
		Display string `json:"display"`
	} `json:"groups"`
	IgnoreLabels []string `json:"ignoreLabels"`
}

// Load processes command-line arguments to construct a Config object.
func Load(args Args, token string) (*Config, error) {
	conf, err := loadConfigFile(args.ConfPath)
	if err != nil {
		return nil, err
	}

	conf.Template, err = parseTemplate(args.TemplatePath)
	if err != nil {
		return nil, err
	}

	conf.FilterLabel = args.FilterLabel

	conf.GitHubToken = strings.TrimSpace(token)
	if conf.GitHubToken == "" {
		return nil, errNoToken
	}

	return conf, nil
}

func loadConfigFile(path string) (*Config, error) {
	confBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	var confFile configFile
	if err := yaml.Unmarshal(confBytes, &confFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from %s: %w", path, err)
	}

	conf := Config{
		Repository:   confFile.Repository,
		GroupLabels:  make(map[string]string, len(confFile.Groups)+1),
		GroupOrder:   make([]string, len(confFile.Groups)+1),
		IgnoreLabels: make(map[string]struct{}, len(confFile.IgnoreLabels)),
	}

	for i, group := range confFile.Groups {
		conf.GroupOrder[i] = group.Label
		conf.GroupLabels[group.Label] = group.Display
	}

	// add built-in defaults for PRs with no other labels
	conf.GroupOrder[len(conf.GroupOrder)-1] = NoGroupKey
	conf.GroupLabels[NoGroupKey] = noGroupDisplay

	for _, lbl := range confFile.IgnoreLabels {
		conf.IgnoreLabels[lbl] = struct{}{}
	}

	return &conf, nil
}

func parseTemplate(path string) (*template.Template, error) {
	funcs := template.FuncMap{
		"id": func(s string) string {
			return strings.TrimPrefix(s, ">")
		},
	}

	return template.New(filepath.Base(path)).Funcs(funcs).ParseFiles(path)
}
