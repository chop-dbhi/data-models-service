package client

import "testing"

func TestClient(t *testing.T) {
	c, err := New(DefaultServiceURL)
	if err != nil {
		t.Fatalf("Error initializing client: %s", err)
	}

	models, err := c.Models()
	if err != nil {
		t.Fatalf("Error fetching models: %s", err)
	}

	if models.Len() == 0 {
		t.Skip("No models for remaining tests")
	}

	model := models.List()[0]

	if _, err = c.ModelRevisions(model.Name); err != nil {
		t.Logf("Error fetching model revisions: %s", err)
	}

	if _, err = c.ModelRevision(model.Name, model.Version); err != nil {
		t.Logf("Error fetching model revision: %s", err)
	}

	if _, err = c.Schema(model.Name, model.Version); err != nil {
		t.Logf("Error fetching schema: %s", err)
	}
}
