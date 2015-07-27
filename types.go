package main

type FileType int

func (f FileType) String() string {
	return fileTypeStrings[f]
}

const (
	UnknownType FileType = iota
	FieldsFile
	TablesFile
	SchemataFile
	ReferencesFile
	IndexesFile
	ConstraintsFile
	MappingsFile
	ModelsFile
)

var fileTypeStrings = map[FileType]string{
	UnknownType:     "unknown",
	FieldsFile:      "fields",
	TablesFile:      "tables",
	SchemataFile:    "schema",
	IndexesFile:     "indexes",
	ConstraintsFile: "constraints",
	MappingsFile:    "mappings",
	ModelsFile:      "models",
}

// Mapping of file types to their minimum required fields.
var FileTypeFields = map[FileType][]string{
	FieldsFile: {
		"model",
		"version",
		"table",
		"field",
		"description",
	},

	TablesFile: {
		"model",
		"version",
		"table",
		"description",
	},

	SchemataFile: {
		"model",
		"version",
		"table",
		"field",
		"type",
		"length",
		"precision",
		"scale",
		"default",
	},

	ConstraintsFile: {
		"model",
		"version",
		"table",
		"field",
		"type",
		"name",
	},

	IndexesFile: {
		"model",
		"version",
		"table",
		"field",
		"name",
		"order",
	},

	ReferencesFile: {
		"version",
		"table",
		"field",
		"ref_table",
		"ref_field",
		"name",
	},

	MappingsFile: {
		"source_model",
		"source_version",
		"source_table",
		"source_field",
		"target_model",
		"target_version",
		"target_table",
		"target_field",
		"comment",
	},

	ModelsFile: {
		"model",
		"version",
		"label",
		"description",
		"url",
	},
}

// Explict order since the tables file is a subset of fields.
// TODO(bjr): change table fields to not be ambiguous
var fileTypesOrder = []FileType{
	FieldsFile,
	SchemataFile,
	IndexesFile,
	ConstraintsFile,
	ReferencesFile,
	TablesFile,
	MappingsFile,
	ModelsFile,
}

func hasFields(header, fields []string) bool {
	m := make(map[string]struct{}, len(header))

	for _, f := range header {
		m[f] = struct{}{}
	}

	for _, f := range fields {
		if _, ok := m[f]; !ok {
			return false
		}
	}

	return true
}

// detectFileType takes a header and attempts to detect the file type based
// on the fields.
func detectFileType(header []string) FileType {
	for _, fileType := range fileTypesOrder {
		if hasFields(header, FileTypeFields[fileType]) {
			return fileType
		}
	}

	return UnknownType
}
