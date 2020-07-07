// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	args := Args{
		ConfPath:     "testdata/config.yaml",
		TemplatePath: "testdata/template.tpl",
		FilterLabel:  "v1.2.0",
	}

	conf, err := Load(args, "xxxyyy")

	require.NoError(t, err)
	require.Equal(t, "elastic/cloud-on-k8s", conf.Repository)
	require.Equal(t, "v1.2.0", conf.FilterLabel)
	require.Equal(t, "xxxyyy", conf.GitHubToken)
	require.NotNil(t, conf.Template)

	wantGroupLabels := map[string]string{
		">breaking":    "Breaking changes",
		">deprecation": "Deprecations",
		">feature":     "New features",
		">enhancement": "Enhancements",
		">bug":         "Bug fixes",
	}

	require.Equal(t, len(wantGroupLabels), len(conf.GroupLabels), "Size of groupLabels do not match")

	for k, v := range wantGroupLabels {
		require.Equal(t, v, conf.GroupLabels[k], "%s in groupLabels do not match", k)
	}

	wantGroupOrder := []string{
		">breaking",
		">deprecation",
		">feature",
		">enhancement",
		">bug",
	}

	require.Equal(t, wantGroupOrder, conf.GroupOrder)

	wantIgnoreLabels := map[string]struct{}{
		">non-issue":                 {},
		">refactoring":               {},
		">docs":                      {},
		">test":                      {},
		":ci":                        {},
		"backport":                   {},
		"exclude-from-release-notes": {},
	}

	require.Equal(t, len(wantIgnoreLabels), len(conf.IgnoreLabels), "Size of ignoreLabels do not match")

	for k := range wantIgnoreLabels {
		_, found := conf.IgnoreLabels[k]
		require.True(t, found, "%s not in ignoreLabels", k)
	}
}
