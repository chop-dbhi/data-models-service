package client

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Attrs is a map of string key/value pairs.
type Attrs map[string]string

type Release struct {
	Level  string `json:"level"`
	Serial string `json:"serial"`
}

type Model struct {
	Label       string `json:"label"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	URL         string `json:"url"`

	Release *Release `json:"release"`
	Tables  *Tables  `json:"tables"`
	Schema  *Schema  `json:"-"`

	Path string `json:"-"`
}

func (m *Model) String() string {
	if m.Label != "" {
		return m.Label
	}

	return fmt.Sprintf("%s/%s", m.Name, m.Version)
}

func (m *Model) URLPath() string {
	return fmt.Sprintf("%s/%s", m.Name, m.Version)
}

func (m *Model) URLSlug() string {
	return fmt.Sprintf("%s-%s", m.Name, m.Version)
}

type Table struct {
	Name        string  `json:"name"`
	Label       string  `json:"label"`
	Description string  `json:"description"`
	Fields      *Fields `json:"fields"`

	Model *Model `json:"-"`
	Attrs Attrs  `json:"-"`
}

func (t *Table) String() string {
	if t.Label != "" {
		return t.Label
	}

	return t.Name
}

func (t *Table) URLPath() string {
	return fmt.Sprintf("%s/%s", t.Model.URLPath(), t.Name)
}

func (t *Table) URLSlug() string {
	return fmt.Sprintf("%s-%s", t.Model.URLSlug(), t.Name)
}

type Field struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Required    bool   `json:"required"`

	// Schema fields
	Type      string `json:"type"`
	Length    int    `json:"length"`
	Precision int    `json:"precision"`
	Scale     int    `json:"scale"`
	Default   string `json:"default"`

	Table *Table `json:"-"`

	Mappings []*Mapping `json:"-"`

	// The field this was renamed from in the previous version.
	RenamedFrom *Field `json:"-"`

	// The field this was renamed to in the next version.
	RenamedTo *Field `json:"-"`

	// The field this field references.
	References *Reference `json:"-"`

	// Fields that reference this field.
	InboundRefs []*Reference `json:"-"`

	Attrs Attrs `json:"-"`
}

func (f *Field) String() string {
	if f.Label != "" {
		return f.Label
	}

	return f.Name
}

func (f *Field) URLPath() string {
	return fmt.Sprintf("%s/%s", f.Table.URLPath(), f.Name)
}

func (f *Field) URLSlug() string {
	return fmt.Sprintf("%s-%s", f.Table.URLSlug(), f.Name)
}

// Models is a sortable slice of models by name then semantic version.
type Models struct {
	m map[string]map[string]*Model
	l []*Model
}

func (ms *Models) Len() int {
	return len(ms.m)
}

// TODO: handle version numbers correctly.
func (ms *Models) Less(i, j int) bool {
	a := ms.l[i]
	b := ms.l[j]

	if a.Name < b.Name {
		return true
	} else if a.Name > b.Name {
		return false
	}

	return a.Version < b.Version
}

func (ms *Models) Swap(i, j int) {
	ms.l[i], ms.l[j] = ms.l[j], ms.l[i]
}

func (ms *Models) Keys() []string {
	keys := make([]string, len(ms.l))

	for i, m := range ms.l {
		keys[i] = m.Name
	}

	return keys
}

func (ms *Models) Add(m *Model) {
	n := strings.ToLower(m.Name)
	v := strings.ToLower(m.Version)

	if ms.m == nil {
		ms.m = make(map[string]map[string]*Model)
	}

	if ix, ok := ms.m[n]; !ok {
		ms.m[n] = map[string]*Model{v: m}
		ms.l = append(ms.l, m)
		sort.Sort(ms)
	} else if _, ok := ix[v]; !ok {
		ix[v] = m
		ms.l = append(ms.l, m)
		sort.Sort(ms)
	}
}

func (ms *Models) Get(n, v string) *Model {
	n = strings.ToLower(n)
	v = strings.ToLower(v)

	if ix, ok := ms.m[n]; ok {
		return ix[v]
	}

	return nil
}

func (ms *Models) Versions(n string) []*Model {
	n = strings.ToLower(n)

	if _, ok := ms.m[n]; !ok {
		return nil
	}

	i := 0
	models := make([]*Model, len(ms.m[n]))

	for _, m := range ms.m[n] {
		models[i] = m
		i++
	}

	return models
}

func (ms *Models) Latest() *Model {
	if ms.Len() == 0 {
		return nil
	}

	return ms.l[len(ms.l)-1]
}

func (ms *Models) List() []*Model {
	return ms.l
}

// MarshalJSON implemenfs the json.Marshaler interface. The marshaled value
// is a sorted list of tables.
func (ms *Models) MarshalJSON() ([]byte, error) {
	return json.Marshal(ms.l)
}

func (ms *Models) UnmarshalJSON(b []byte) error {
	var aux []*Model

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	if ms.m == nil {
		ms.m = make(map[string]map[string]*Model)
	}

	for _, m := range aux {
		ms.Add(m)
	}

	return nil
}

// Tables is a set of tables.
type Tables struct {
	m map[string]*Table
	l []*Table
}

// List returns a slice of tables.
func (ts *Tables) List() []*Table {
	return ts.l
}

// Names returns a sorted list of table names.
func (ts *Tables) Names() []string {
	l := make([]string, len(ts.l))

	for i, t := range ts.l {
		l[i] = t.Name
	}

	return l
}

// Add adds a table to the set.
func (ts *Tables) Add(t *Table) {
	k := strings.ToLower(t.Name)

	if ts.m == nil {
		ts.m = make(map[string]*Table)
	}

	if _, ok := ts.m[k]; !ok {
		ts.m[k] = t
		ts.l = append(ts.l, t)
		sort.Sort(ts)
	}
}

// Get returns a table by name.
func (ts *Tables) Get(n string) *Table {
	return ts.m[strings.ToLower(n)]
}

func (ts *Tables) Len() int {
	return len(ts.m)
}

func (ts *Tables) Less(i, j int) bool {
	a := ts.l[i]
	b := ts.l[j]

	return a.Name < b.Name
}

func (ts *Tables) Swap(i, j int) {
	ts.l[i], ts.l[j] = ts.l[j], ts.l[i]
}

// MarshalJSON implements the json.Marshaler interface. The marshaled value
// is a sorted list of tables.
func (ts *Tables) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.l)
}

// UnmarshalJSON unmarshals the bytes into the set.
func (ts *Tables) UnmarshalJSON(b []byte) error {
	var aux []*Table

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	if ts.m == nil {
		ts.m = make(map[string]*Table)
	}

	for _, t := range aux {
		ts.Add(t)
	}

	return nil
}

// Fields is a set of fields.
type Fields struct {
	m map[string]*Field
	l []*Field
}

func (fs *Fields) List() []*Field {
	return fs.l
}

// Names returns a sorted list of field names.
func (fs *Fields) Names() []string {
	l := make([]string, len(fs.l))

	for i, t := range fs.l {
		l[i] = t.Name
	}

	return l
}

// Add adds a field to the set.
func (fs *Fields) Add(t *Field) {
	k := strings.ToLower(t.Name)

	if fs.m == nil {
		fs.m = make(map[string]*Field)
	}

	if _, ok := fs.m[k]; !ok {
		fs.m[k] = t
		fs.l = append(fs.l, t)
		sort.Sort(fs)
	}
}

// Get gets a field by name.
func (fs *Fields) Get(n string) *Field {
	return fs.m[strings.ToLower(n)]
}

func (fs *Fields) Len() int {
	return len(fs.m)
}

func (fs *Fields) Less(i, j int) bool {
	a := fs.l[i]
	b := fs.l[j]

	return a.Name < b.Name
}

func (fs *Fields) Swap(i, j int) {
	fs.l[i], fs.l[j] = fs.l[j], fs.l[i]
}

// MarshalJSON implemenfs the json.Marshaler interface. The marshaled value
// is a sorted list of tables.
func (fs *Fields) MarshalJSON() ([]byte, error) {
	return json.Marshal(fs.l)
}

func (fs *Fields) UnmarshalJSON(b []byte) error {
	var aux []*Field

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	if fs.m == nil {
		fs.m = make(map[string]*Field)
	}

	for _, f := range aux {
		fs.Add(f)
	}

	return nil
}
