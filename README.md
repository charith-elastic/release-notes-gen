Release notes generator
=======================

This tool generates release notes from all PRs labeled with a specific label. If any issues are linked to a PR, they will be included in the output as well.


Prerequisites
--------------

Create a GitHub token by going to https://github.com/settings/tokens. The token must have `repo:status` and `public_repo` scopes. Enable SSO for the token as well.


Usage
-----

```
Usage:
  release-notes-gen [flags]

Examples:
GITHUB_TOKEN=xxxyyy ./release-notes-gen --conf=example/config.yaml --template=example/template.tpl --label=v1.2.0

Flags:
      --conf string       Path to the configuration file
  -h, --help              help for release-notes-gen
      --label string      Label to filter PRs by (e.g. v1.2.0)
      --template string   Path to the release notes template
  -v, --version           version for release-notes-gen
```


### Configuration file format

```yaml
# Name of the repository to inspect
repository: elastic/cloud-on-k8s
# Labels used for grouping the pull requests into categories. If a PR has multiple labels, the label appearing
# earliest in the list wins.
groups:
  - label: ">breaking"
    display: "Breaking changes"
  - label: ">deprecation"
    display: "Deprecations"
  - label: ">feature"
    display: "New features"
  - label: ">enhancement"
    display: "Enhancements"
  - label: ">bug"
    display: "Bug fixes"
# Ignore any PRs with the following labels
ignoreLabels:
  - ">non-issue"
  - ">refactoring"
  - ">docs"
  - ">test"
  - ":ci"
  - "backport"
  - "exclude-from-release-notes"
```
