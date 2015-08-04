package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var (
	defaultFormat = "html"

	mimetypes = map[string]string{
		"text/markdown":    "markdown",
		"text/html":        "html",
		"application/json": "json",
	}

	queryFormats = map[string]string{
		"md":       "markdown",
		"markdown": "markdown",
		"html":     "html",
		"json":     "json",
	}

	userAgent = "DataModelsService/%s (+https://github.com/chop-dbhi/data-models-service)"
)

func jsonResponse(w http.ResponseWriter, d interface{}) error {
	e := json.NewEncoder(w)

	if err := e.Encode(d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}

// detectFormat applies content negotiation logic to determine the
// appropriate response representation.
func detectFormat(w http.ResponseWriter, r *http.Request) string {
	var (
		ok     bool
		format string
	)

	format = queryFormats[strings.ToLower(r.URL.Query().Get("format"))]

	// Query parameter
	if format == "" {
		// Accept header
		acceptType := r.Header.Get("Accept")
		acceptType, _, _ = mime.ParseMediaType(acceptType)

		// Fallback to default
		if format, ok = mimetypes[acceptType]; !ok {
			format = defaultFormat
		}
	}

	var contentType string

	switch format {
	case "html":
		contentType = "text/html"
	case "markdown":
		contentType = "text/markdown"
	case "json":
		contentType = "application/json"
	}

	w.Header().Set("user-agent", fmt.Sprintf(userAgent, progVersion))
	w.Header().Set("content-type", contentType)

	return format
}

func httpIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	switch detectFormat(w, r) {
	case "html":
		RenderIndexHTML(w)
	case "json":
		jsonResponse(w, map[string]interface{}{
			"name":    serviceName,
			"version": progVersion,
		})
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func httpModels(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := map[string]interface{}{
		"Title": "Models",
		"Items": dataModelCache.List(),
	}

	switch detectFormat(w, r) {
	case "markdown":
		RenderModelsMarkdown(w, data)
	case "html":
		RenderModelsHTML(w, data)
	case "json":
		jsonResponse(w, data["Items"])
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func httpModel(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")

	m := dataModelCache.Versions(n)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": m[0].Name,
		"Items": m,
	}

	switch detectFormat(w, r) {
	case "markdown":
		RenderModelMarkdown(w, data)
	case "html":
		RenderModelHTML(w, data)
	case "json":
		jsonResponse(w, m)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func httpModelVersion(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")

	m := dataModelCache.Get(n, v)

	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch detectFormat(w, r) {
	case "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderModelVersionMarkdown(w, m)
	case "html":
		w.Header().Set("content-type", "text/html")
		RenderModelVersionHTML(w, m)
	case "json":
		jsonResponse(w, m)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func httpTable(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")
	tn := p.ByName("table")

	var (
		m *Model
		t *Table
	)

	if m = dataModelCache.Get(n, v); m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if t = m.Tables.Get(tn); t == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch detectFormat(w, r) {
	case "json":
		jsonResponse(w, t)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func httpField(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")
	tn := p.ByName("table")
	fn := p.ByName("field")

	var (
		m *Model
		t *Table
		f *Field
	)

	if m = dataModelCache.Get(n, v); m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if t = m.Tables.Get(tn); t == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if f = t.Fields.Get(fn); f == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch detectFormat(w, r) {
	case "json":
		jsonResponse(w, f)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func httpCompareModels(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n1 := p.ByName("name1")
	v1 := p.ByName("version1")

	m1 := dataModelCache.Get(n1, v1)

	if m1 == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	n2 := p.ByName("name2")
	v2 := p.ByName("version2")

	m2 := dataModelCache.Get(n2, v2)

	if m2 == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch detectFormat(w, r) {
	case "md", "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderModelCompareMarkdown(w, m1, m2)
	case "", "html":
		w.Header().Set("content-type", "text/html")
		RenderModelCompareHTML(w, m1, m2)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func verifyGithubSignature(sig string, r io.Reader) bool {
	mac := hmac.New(sha1.New, []byte(secret))
	io.Copy(mac, r)
	expected := fmt.Sprintf("sha1=%x", mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(sig))
}

// TODO: does each repo have a webhook?
func httpUpdateRepos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// If no secret has been supplied, update at will.
	if secret == "" {
		updateRepos()
		return
	}

	// Check for Github's webhook signature.
	if sig := r.Header.Get("X-Hub-Signature"); sig != "" {
		defer r.Body.Close()

		if verifyGithubSignature(sig, r.Body) {
			updateRepos()
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func httpModelSchema(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	n := p.ByName("name")
	v := p.ByName("version")

	var (
		m   *Model
		err error
	)

	if m = dataModelCache.Get(n, v); m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	aux := make(map[string]interface{})

	aux["model"] = m.Name
	aux["version"] = m.Version
	aux["tables"] = m.Tables
	aux["schema"] = m.schema

	switch detectFormat(w, r) {
	case "json":
		if err = jsonResponse(w, aux); err != nil {
			w.Write([]byte(fmt.Sprintf("erroring marshaling %s schema: %s", m, err)))
		}
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func httpReposList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	switch detectFormat(w, r) {
	case "md", "markdown":
		w.Header().Set("content-type", "text/markdown")
		RenderReposMarkdown(w, registeredRepos)
	case "", "html":
		w.Header().Set("content-type", "text/html")
		RenderReposHTML(w, registeredRepos)
	case "json":
		jsonResponse(w, registeredRepos)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}
