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
