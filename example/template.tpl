:issue: https://github.com/{{.Repo}}/issues/
:pull: https://github.com/{{.Repo}}/pull/

[[release-notes-{{.FilterLabel}}]]
== {n} version {{.FilterLabel}}
{{range $group := .GroupOrder -}}
{{with (index $.Groups $group)}}
[[{{- id $group -}}-{{$.FilterLabel}}]]
[float]
=== {{index $.GroupLabels $group}}
{{range .}}
* {{.Title}} {pull}{{.Number}}[#{{.Number}}]{{with .Issues -}}
{{$length := len .}} (issue{{if gt $length 1}}s{{end}}: {{range $idx, $el := .}}{{if $idx}}, {{end}}{issue}{{$el}}[#{{$el}}]{{end}})
{{- end}}
{{- end}}
{{- end}}
{{end}}
