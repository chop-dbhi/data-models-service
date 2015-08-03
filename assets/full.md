# {{.}}

{{if .Description}}{{.Description}}{{end}}

- Version: {{.Version}}{{if .Release.Level}}
- Release: {{.Release.Level}}+{{.Release.Serial}}{{end}}
- URL: {{.URL}}

## Tables

{{range .Tables.List}}- [{{.}}](#{{.URLSlug}})
{{end}}

{{range .Tables.List}}## {{.}} {#{{.URLSlug}}}

{{.Description}}

**Fields**

{{range .Fields.List}}- [{{.}}](#{{.URLSlug}})
{{end}}
{{range .Fields.List}}#### {{.}} {#{{.URLSlug}}}

{{if .References}}*Refers to: [{{.References.Field.Table}}](#{{.References.Field.Table.URLSlug}}) / [{{.References.Field}}](#{{.References.Field.URLSlug}})*{{end}}

{{.Description}}

{{if .Type}}##### Schema

- Type: `{{.Type}}`{{if .Length}}
- Length: {{.Length}}{{end}}{{if .Precision}}
- Precision: {{.Precision}}{{end}}{{if .Scale}}
- Scale: {{.Scale}}{{end}}
{{end}}

{{if .Mappings}}##### Mappings

Model | Table | Field | Comment
------|-------|-------|--------
{{range .Mappings}}[{{.Field.Table.Model}}](/models/{{.Field.Table.Model.URLPath}}) | [{{.Field.Table}}](/models/{{.Field.Table.Model.URLPath}}#{{.Field.Table.URLSlug}}) | [{{.Field}}](/models/{{.Field.Table.Model.URLPath}}#{{.Field.URLSlug}}) | {{.Comment}}
{{end}}
{{end}}

{{if .InboundRefs}}##### Inbound References

*Total: {{len .InboundRefs}}*

Table | Field | Name
------|-------|-----
{{range .InboundRefs}}[{{.Field.Table}}](#{{.Field.Table.URLSlug}}) | [{{.Field}}](#{{.Field.URLSlug}}) | {{.}}
{{end}}
{{end}}

{{end}}
{{end}}
