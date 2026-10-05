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
	if pageData, ok := data.(PageData); ok {
		if pageData.Language == "" {
			pageData.Language = languageEnglish
		}
		if pageData.Text.LandingTitle == "" {
			pageData.Text = translationsFor(pageData.Language)
		}
		data = pageData
	}

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

func renderGuestPage(
	w http.ResponseWriter,
	r *http.Request,
	name string,
	data PageData,
) {
	data.Language = languageForHost(r.Host)
	data.Text = translationsFor(data.Language)

	switch data.ErrorKey {
	case errorInvalidInvitation:
		data.Error = data.Text.InvalidInvitation
	case errorAllergies:
		data.Error = data.Text.ErrorAllergies
	case errorAccommodation:
		data.Error = data.Text.ErrorAccommodation
	case errorMessage:
		data.Error = data.Text.ErrorMessage
	case errorSave:
		data.Error = data.Text.ErrorSave
	}

	render(w, name, data)
}
