{
    "private": true,
    "workspaces": [
    {{- $first := true}}
    {{- range $v := .ServiceList}}
        {{- if $first}}
            {{- $first = false}}
        {{- else}},
        {{- end}}
        "{{$v}}"
    {{- end}}
    ]
}
