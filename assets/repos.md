# Repositories

{{range .}}- {{.URL}}
    - Branch: {{.Branch}}
    - Commit: {{.CommitSHA1}}
    - Commit Date: {{.CommitTime.Local}}
    - Fetched: {{.FetchTime.Local}}
{{end}}
