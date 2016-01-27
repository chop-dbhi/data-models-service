package main

import (
	"encoding/csv"
	"io"
	"strings"

	dms "github.com/chop-dbhi/data-models-service/client"
)

// UniversalReader wraps an io.Reader to replace carriage returns with newlines.
// This is used with the csv.Reader so it can properly delimit lines.
type UniversalReader struct {
	r io.Reader
}

func (r *UniversalReader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)

	// Replace carriage returns with newlines
	for i, b := range buf {
		if b == '\r' {
			buf[i] = '\n'
		}
	}

	return n, err
}

type MapCSVReader struct {
	fields []string

	csv *csv.Reader
}

func (r *MapCSVReader) zip(keys, values []string) dms.Attrs {
	m := make(dms.Attrs, len(keys))

	for i, k := range keys {
		m[strings.TrimSpace(k)] = strings.TrimSpace(values[i])
	}

	return m
}

func (r *MapCSVReader) Fields() []string {
	if r.fields == nil {
		fields, err := r.csv.Read()

		if err != nil {
			return nil
		}

		r.fields = fields
	}

	return r.fields
}

func (r *MapCSVReader) Read() (dms.Attrs, error) {
	// First iteration.
	if r.fields == nil {
		r.Fields()
	}

	values, err := r.csv.Read()

	if err != nil {
		return nil, err
	}

	m := r.zip(r.fields, values)

	return m, err
}

func (r *MapCSVReader) ReadAll() ([]dms.Attrs, error) {
	records := make([]dms.Attrs, 0)

	for {
		record, err := r.Read()

		if record == nil {
			break
		}

		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func NewMapCSVReader(r io.Reader) *MapCSVReader {
	cr := csv.NewReader(&UniversalReader{r})

	cr.LazyQuotes = true
	cr.TrimLeadingSpace = true

	return &MapCSVReader{
		csv: cr,
	}
}
