package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRenderSuccessPageWithNilInvitation(t *testing.T) {
	rr := httptest.NewRecorder()

	render(rr, "success.html", PageData{RSVP: &RSVP{Attending: true}})

	if rr.Code == http.StatusInternalServerError {
		t.Fatalf("success page should render without a nil invitation, got %d: %s", rr.Code, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "RSVP received") {
		t.Fatalf("success page body missing expected content: %s", rr.Body.String())
	}
}

func TestRenderRSVPPageWithoutSavedRSVP(t *testing.T) {
	rr := httptest.NewRecorder()
	invitation := Invitation{Code: "ALICE-BOB-7K2P", CoupleName: "Alice & Bob"}

	render(rr, "rsvp.html", PageData{Invitation: &invitation})

	if rr.Code == http.StatusInternalServerError {
		t.Fatalf("rsvp page should render without a saved RSVP, got %d: %s", rr.Code, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "Send my RSVP") {
		t.Fatalf("rsvp page body missing expected button: %s", rr.Body.String())
	}

	formPosition := strings.Index(rr.Body.String(), `id="rsvp"`)
	informationPosition := strings.Index(rr.Body.String(), `id="weekend"`)
	if formPosition < 0 || informationPosition < 0 || formPosition > informationPosition {
		t.Fatal("rsvp form should appear before the guest information sections")
	}

	for _, expected := range []string{
		"Luky &amp; Lelaina",
		"Dear Alice &amp; Bob,",
		"reserved 2 seats",
		"28 August",
		"Statek &Uacute;jezd u Pl&aacute;nice",
		"field behind the venue",
		"Honzí and Anna",
		"/static/wedding.ics",
	} {
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("rsvp page body missing %q", expected)
		}
	}
}

func TestLandingPageIncludesGuestInformation(t *testing.T) {
	rr := httptest.NewRecorder()
	render(rr, "landing.html", PageData{})

	for _, expected := range []string{
		"28 August",
		"Statek &Uacute;jezd u Pl&aacute;nice",
		"Open in Google Maps",
		"/static/wedding.ics",
		"id=\"rsvp\"",
	} {
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("landing page body missing %q", expected)
		}
	}
}

func TestGuestLanguageFollowsSubdomain(t *testing.T) {
	tests := []struct {
		host     string
		language string
		text     string
	}{
		{host: "cs.example.com", language: "cs", text: "Bereme se"},
		{host: "cz.example.com:8080", language: "cs", text: "Bereme se"},
		{host: "de.example.com", language: "de", text: "Wir heiraten"},
		{host: "en.example.com", language: "en", text: "getting married"},
		{host: "localhost:8080", language: "en", text: "getting married"},
	}

	for _, test := range tests {
		t.Run(test.host, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://"+test.host+"/", nil)
			res := httptest.NewRecorder()

			renderGuestPage(res, req, "landing.html", PageData{})

			if !strings.Contains(res.Body.String(), `lang="`+test.language+`"`) {
				t.Errorf("page language does not match host %q: %s", test.host, res.Body.String())
			}
			if !strings.Contains(res.Body.String(), test.text) {
				t.Errorf("page text for host %q does not contain %q", test.host, test.text)
			}
		})
	}
}

func TestGermanRSVPUsesTranslatedLabelsAndStableValues(t *testing.T) {
	invitation := Invitation{Code: "ALICE-BOB-7K2P", CoupleName: "Alice & Bob"}
	req := httptest.NewRequest(http.MethodGet, "http://de.example.com/rsvp", nil)
	res := httptest.NewRecorder()

	renderGuestPage(res, req, "rsvp.html", PageData{Invitation: &invitation})
	body := res.Body.String()

	for _, expected := range []string{
		`lang="de"`,
		"Allergien oder Ernährungswünsche",
		"Rückmeldung senden",
		`name="accommodation" value="I will organise myself"`,
	} {
		if !strings.Contains(body, expected) {
			t.Errorf("German RSVP page missing %q", expected)
		}
	}
}

func TestInvalidInvitationErrorUsesHostLanguage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://cs.example.com/", nil)
	res := httptest.NewRecorder()

	renderGuestPage(res, req, "landing.html", PageData{ErrorKey: errorInvalidInvitation})

	if !strings.Contains(res.Body.String(), "Tento kód pozvánky se nepodařilo najít") {
		t.Fatalf("Czech page should show a Czech invitation error: %s", res.Body.String())
	}
}

func TestSubmitRSVPStoresPerGuestAttendanceAndAllergies(t *testing.T) {
	form := url.Values{}
	form.Set("code", "ALICE-BOB-7K2P")
	form.Set("guest_0_attending", "Alice")
	form.Set("guest_0_allergies", "Peanuts")
	form.Set("guest_1_attending", "Bob")
	form.Set("guest_1_allergies", "Gluten")
	form.Set("accommodation", "I will organise myself")
	form.Set("message", "Thanks!")

	req := httptest.NewRequest(http.MethodPost, "/rsvp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()
	submitRSVP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("submitRSVP should succeed, got %d: %s", res.Code, res.Body.String())
	}

	storeMu.RLock()
	defer storeMu.RUnlock()

	stored, ok := rsvpStore["ALICE-BOB-7K2P"]
	if !ok {
		t.Fatal("RSVP was not saved")
	}

	if len(stored.Guests) != 2 {
		t.Fatalf("expected 2 guest records, got %d", len(stored.Guests))
	}

	if !stored.Guests[0].Attending || stored.Guests[0].Allergies != "Peanuts" {
		t.Fatalf("first guest record was not saved correctly: %+v", stored.Guests[0])
	}

	if !stored.Guests[1].Attending || stored.Guests[1].Allergies != "Gluten" {
		t.Fatalf("second guest record was not saved correctly: %+v", stored.Guests[1])
	}

	if stored.Accommodation != "I will organise myself" {
		t.Fatalf("accommodation was not saved correctly: %q", stored.Accommodation)
	}
}
