package main

import (
	"net/http"
	"strings"
	"unicode/utf8"
)

func calendarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := normalizeCode(r.URL.Query().Get("code"))
	if code == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if _, ok := getInvitation(code); !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	language := languageForHost(r.Host)
	content := weddingCalendar(language)
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="wedding.ics"`)
	w.Header().Set("Cache-Control", "private, no-store")
	if r.Method == http.MethodHead {
		return
	}

	_, _ = w.Write([]byte(content))
}

func weddingCalendar(language Language) string {
	text := translationsFor(language)
	locale := string(language)
	lines := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:-//Luky and Lelaina//Wedding RSVP//EN",
		"CALSCALE:GREGORIAN",
		"METHOD:PUBLISH",
		"X-WR-CALNAME;LANGUAGE=" + locale + ":" + escapeCalendarText(text.CalendarTitle),
		"BEGIN:VEVENT",
		"UID:wedding-20270828@luky-lelaina",
		"DTSTAMP:20261001T000000Z",
		"DTSTART;VALUE=DATE:20270828",
		"DTEND;VALUE=DATE:20270829",
		"SUMMARY;LANGUAGE=" + locale + ":" + escapeCalendarText(text.CalendarTitle),
		"LOCATION:Statek Újezd u Plánice",
		"DESCRIPTION;LANGUAGE=" + locale + ":" + escapeCalendarText(text.CalendarDescription),
		"END:VEVENT",
		"END:VCALENDAR",
	}

	var calendar strings.Builder
	for _, line := range lines {
		calendar.WriteString(foldCalendarLine(line))
		calendar.WriteString("\r\n")
	}
	return calendar.String()
}

func escapeCalendarText(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	return strings.NewReplacer(
		"\\", "\\\\",
		"\n", "\\n",
		",", "\\,",
		";", "\\;",
	).Replace(value)
}

func foldCalendarLine(line string) string {
	var folded strings.Builder
	lineLength := 0
	for _, character := range line {
		characterLength := utf8.RuneLen(character)
		if lineLength+characterLength > 75 {
			folded.WriteString("\r\n ")
			lineLength = 1
		}
		folded.WriteRune(character)
		lineLength += characterLength
	}
	return folded.String()
}
