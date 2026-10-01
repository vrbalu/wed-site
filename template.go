package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

// tmpl is the shared HTML template set. The function map is intentionally small,
// but it is the extension point for adding helpers such as formatting dates,
// resolving names, or enriching template output without scattering logic inside
// the templates themselves.
var tmpl = template.Must(
	template.New("").Funcs(template.FuncMap{
		"invitationGuestNames": invitationGuestNames,
	}).ParseGlob("templates/*.html"),
)

func render(
	w http.ResponseWriter,
	name string,
	data any,
) {
	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	var buf bytes.Buffer

	if err := tmpl.ExecuteTemplate(
		&buf,
		name,
		data,
	); err != nil {
		log.Printf(
			"template %s: %v",
			name,
			err,
		)
		http.Error(
			w,
			"Internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Printf(
			"write template %s: %v",
			name,
			err,
		)
	}
}
