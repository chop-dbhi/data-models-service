package client

import (
	"encoding/json"
	"strings"
)

// A mapping defined a correspondence between two fields. A mapping points
// points to the opposing field and a *comment* which describes the nuances
// of the relationship between the two fields.
type Mapping struct {
	Field   *Field
	Comment string
}

// Reference declares that the source field is a reference to the target field.
type Reference struct {
	Name  string
	Field *Field

	Attrs Attrs `json:"-"`
}

func (r *Reference) String() string {
	return r.Name
}

func (r *Reference) MarshalJSON() ([]byte, error) {
	aux := map[string]string{
		"name":  r.Name,
		"table": r.Field.Table.Name,
		"field": r.Field.Name,
	}

	return json.Marshal(aux)
}

// Schema contains constraints and indexes for a model.
type Schema struct {
	// Schematic components.
	PrimaryKeys  map[string]*PrimaryKey
	ForeignKeys  []*ForeignKey
	NotNullables []*NotNullable
	Uniques      map[string]*Unique
	Indexes      map[string]*Index
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	pks := make([]*PrimaryKey, len(s.PrimaryKeys))
	uniqs := make([]*Unique, len(s.Uniques))
	indexes := make([]*Index, len(s.Indexes))

	i := 0

	for _, pk := range s.PrimaryKeys {
		pks[i] = pk
		i++
	}

	i = 0

	for _, un := range s.Uniques {
		uniqs[i] = un
		i++
	}

	i = 0

	for _, idx := range s.Indexes {
		indexes[i] = idx
		i++
	}

	aux := map[string]interface{}{
		"indexes": indexes,
		"constraints": map[string]interface{}{
			"foreign_keys": s.ForeignKeys,
			"primary_keys": pks,
			"uniques":      uniqs,
			"not_null":     s.NotNullables,
		},
	}

	return json.Marshal(aux)
}

func (s *Schema) UnmarshalJSON(b []byte) error {
	var (
		err error
		aux map[string]json.RawMessage

		indexes []*Index
	)

	if err = json.Unmarshal(b, &aux); err != nil {
		return err
	}

	if err = json.Unmarshal(aux["indexes"], &indexes); err != nil {
		return err
	}

	s.Indexes = make(map[string]*Index, len(indexes))

	for _, idx := range indexes {
		s.Indexes[idx.Name] = idx
	}

	var constrs map[string]json.RawMessage

	if err = json.Unmarshal(aux["constraints"], &constrs); err != nil {
		return err
	}

	if err = json.Unmarshal(constrs["foreign_keys"], &s.ForeignKeys); err != nil {
		return err
	}

	if err = json.Unmarshal(constrs["not_null"], &s.NotNullables); err != nil {
		return err
	}

	var (
		pks   []*PrimaryKey
		uniqs []*Unique
	)

	if err = json.Unmarshal(constrs["primary_keys"], &pks); err != nil {
		return err
	}

	s.PrimaryKeys = make(map[string]*PrimaryKey, len(pks))

	for _, pk := range pks {
		s.PrimaryKeys[pk.Name] = pk
	}

	if err = json.Unmarshal(constrs["uniques"], &uniqs); err != nil {
		return err
	}

	s.Uniques = make(map[string]*Unique, len(uniqs))

	for _, uniq := range uniqs {
		s.Uniques[uniq.Name] = uniq
	}

	return nil
}

func (s *Schema) AddPrimaryKey(a Attrs) {
	if s.PrimaryKeys == nil {
		s.PrimaryKeys = make(map[string]*PrimaryKey)
	}

	n := a["name"]

	if pk, ok := s.PrimaryKeys[n]; !ok {
		s.PrimaryKeys[n] = &PrimaryKey{
			Name:   n,
			Table:  a["table"],
			Fields: []string{a["field"]},
		}
	} else {
		pk.Fields = append(pk.Fields, a["field"])
	}
}

func (s *Schema) AddForeignKey(a Attrs) {
	s.ForeignKeys = append(s.ForeignKeys, &ForeignKey{
		Name:        a["name"],
		SourceTable: a["table"],
		SourceField: a["field"],
		TargetTable: a["ref_table"],
		TargetField: a["ref_field"],
	})
}

func (s *Schema) AddNotNullable(a Attrs) {
	s.NotNullables = append(s.NotNullables, &NotNullable{
		Table: a["table"],
		Field: a["field"],
	})
}

func (s *Schema) AddUnique(a Attrs) {
	if s.Uniques == nil {
		s.Uniques = make(map[string]*Unique)
	}

	n := a["name"]

	if un, ok := s.Uniques[n]; !ok {
		s.Uniques[n] = &Unique{
			Name:   n,
			Table:  a["table"],
			Fields: []string{a["field"]},
		}
	} else {
		un.Fields = append(un.Fields, a["field"])
	}
}

func (s *Schema) AddIndex(a Attrs) {
	if s.Indexes == nil {
		s.Indexes = make(map[string]*Index)
	}

	n := a["name"]

	if idx, ok := s.Indexes[n]; !ok {
		var uniq bool

		switch strings.ToLower(a["unique"]) {
		case "yes", "y", "1":
			uniq = true
		}

		s.Indexes[n] = &Index{
			Name:   n,
			Order:  a["order"],
			Table:  a["table"],
			Unique: uniq,
			Fields: []string{a["field"]},
		}
	} else {
		idx.Fields = append(idx.Fields, a["field"])
	}
}

// PrimaryKey is a constraint which declares the field values uniquely define
// a record in the respective table.
type PrimaryKey struct {
	Name   string   `json:"name"`
	Table  string   `json:"table"`
	Fields []string `json:"fields"`
}

// Unique is a constraint which declares the field values be unique for
// a record in the respective table.
type Unique struct {
	Name   string   `json:"name"`
	Table  string   `json:"table"`
	Fields []string `json:"fields"`
}

// ForeignKey is a constraint which declares the field values are constrained
// to values in the referenced fields.
type ForeignKey struct {
	Name        string `json:"name"`
	SourceTable string `json:"source_table"`
	SourceField string `json:"source_field"`
	TargetTable string `json:"target_table"`
	TargetField string `json:"target_field"`
}

// NotNullable is a constraint which declares the field cannot be a null.
type NotNullable struct {
	Table string `json:"table"`
	Field string `json:"field"`
}

// Index represents a schematic index for one or more fields.
type Index struct {
	Name   string   `json:"name"`
	Unique bool     `json:"unique"`
	Order  string   `json:"order"`
	Table  string   `json:"table"`
	Fields []string `json:"fields"`
}
