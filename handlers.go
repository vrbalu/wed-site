package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
)

// --------------------------------------------------
// Landing page
// --------------------------------------------------

func landingHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	code := normalizeCode(
		r.URL.Query().Get("code"),
	)

	if code == "" {
		renderGuestPage(
			w,
			r,
			"landing.html",
			PageData{},
		)
		return
	}

	invitation, ok := getInvitation(code)

	if !ok {
		renderGuestPage(
			w,
			r,
			"landing.html",
			PageData{
				ErrorKey: errorInvalidInvitation,
			},
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/rsvp?code="+urlQuery(invitation.Code),
		http.StatusSeeOther,
	)
}

// --------------------------------------------------
// RSVP
// --------------------------------------------------

func rsvpHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {
	case http.MethodGet:
		showRSVP(w, r)

	case http.MethodPost:
		submitRSVP(w, r)

	default:
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

// GET /rsvp?code=ALICE-BOB-7K2P
func showRSVP(
	w http.ResponseWriter,
	r *http.Request,
) {
	code := normalizeCode(
		r.URL.Query().Get("code"),
	)

	invitation, ok := getInvitation(code)

	if !ok {
		http.Redirect(
			w,
			r,
			"/",
			http.StatusSeeOther,
		)
		return
	}

	rsvp, hasRSVP := repo.GetRSVP(invitation.Code)

	var rsvpPtr *RSVP

	if hasRSVP {
		copy := rsvp
		rsvpPtr = &copy
	}

	renderGuestPage(
		w,
		r,
		"rsvp.html",
		PageData{
			Invitation: &invitation,
			RSVP:       rsvpPtr,
			HasRSVP:    hasRSVP,
		},
	)
}

// POST /rsvp
func submitRSVP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			"Invalid form submission",
			http.StatusBadRequest,
		)
		return
	}

	code := normalizeCode(
		r.FormValue("code"),
	)

	invitation, ok := getInvitation(code)

	if !ok {
		http.Error(
			w,
			"Invalid invitation code",
			http.StatusBadRequest,
		)
		return
	}

	// -----------------------------------------------
	// Attendance by guest name
	// -----------------------------------------------

	guestNames := invitationGuestNames(invitation)
	guests := make([]Guest, 0, len(guestNames))
	allergiesParts := make([]string, 0, len(guestNames))
	anyAttending := false
	guestCount := 0

	for i, guestName := range guestNames {
		attendingValue := r.FormValue(fmt.Sprintf("guest_%d_attending", i))
		attending := attendingValue == guestName
		allergies := strings.TrimSpace(
			r.FormValue(fmt.Sprintf("guest_%d_allergies", i)),
		)

		if len(allergies) > 500 {
			renderGuestPage(
				w,
				r,
				"rsvp.html",
				PageData{
					Invitation: &invitation,
					ErrorKey:   errorAllergies,
				},
			)
			return
		}

		if attending {
			anyAttending = true
			guestCount++
		}

		if allergies != "" {
			allergiesParts = append(allergiesParts, fmt.Sprintf("%s: %s", guestName, allergies))
		}

		guests = append(guests, Guest{
			Name:      guestName,
			Attending: attending,
			Allergies: allergies,
		})
	}

	// -----------------------------------------------
	// Other fields
	// -----------------------------------------------

	accommodation := strings.TrimSpace(
		r.FormValue("accommodation"),
	)
	if accommodation != "I will organise myself" &&
		accommodation != "If you find me place, it will be awsome!" {
		renderGuestPage(
			w,
			r,
			"rsvp.html",
			PageData{
				Invitation: &invitation,
				ErrorKey:   errorAccommodation,
			},
		)
		return
	}

	message := strings.TrimSpace(
		r.FormValue("message"),
	)

	if len(message) > 1000 {
		renderGuestPage(
			w,
			r,
			"rsvp.html",
			PageData{
				Invitation: &invitation,
				ErrorKey:   errorMessage,
			},
		)
		return
	}

	// -----------------------------------------------
	// Save RSVP
	// -----------------------------------------------

	rsvp := RSVP{
		Code:          invitation.Code,
		CoupleName:    invitation.CoupleName,
		Attending:     anyAttending,
		GuestCount:    guestCount,
		Accommodation: accommodation,
		Allergies:     strings.Join(allergiesParts, "; "),
		Message:       message,
		SubmittedAt:   time.Now(),
		Guests:        guests,
	}

	if err := repo.SaveRSVP(rsvp); err != nil {
		renderGuestPage(
			w,
			r,
			"rsvp.html",
			PageData{
				Invitation: &invitation,
				ErrorKey:   errorSave,
			},
		)
		return
	}

	storeMu.Lock()
	rsvpStore[invitation.Code] = rsvp
	storeMu.Unlock()

	renderGuestPage(
		w,
		r,
		"success.html",
		PageData{
			Invitation: &invitation,
			RSVP:       &rsvp,
			HasRSVP:    true,
		},
	)
}

// --------------------------------------------------
// Admin login
// --------------------------------------------------

func adminLoginHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {

	case http.MethodGet:
		render(
			w,
			"admin-login.html",
			PageData{},
		)

	case http.MethodPost:
		password := os.Getenv("ADMIN_PASSWORD")

		if password == "" {
			password = "change-me"
		}

		if r.FormValue("password") != password {
			render(
				w,
				"admin-login.html",
				PageData{
					Error: "Incorrect password.",
				},
			)
			return
		}

		http.SetCookie(
			w,
			&http.Cookie{
				Name:     "wedding_admin",
				Value:    "authenticated",
				Path:     "/admin",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Secure:   r.TLS != nil,
				MaxAge:   8 * 60 * 60,
			},
		)

		http.Redirect(
			w,
			r,
			"/admin",
			http.StatusSeeOther,
		)

	default:
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

// --------------------------------------------------
// Admin dashboard
// --------------------------------------------------

func adminHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if !isAdmin(r) {
		http.Redirect(
			w,
			r,
			"/admin/login",
			http.StatusSeeOther,
		)

		return
	}

	invitations = repo.ListInvitations()
	data := AdminData{
		TotalInvitations: len(invitations),
	}

	for _, invitation := range invitations {
		row := AdminRow{
			Invitation: invitation,
		}

		if rsvp, ok := repo.GetRSVP(invitation.Code); ok {
			data.Responses++
			copy := rsvp
			row.RSVP = &copy

			if rsvp.Attending {
				data.AttendingCouples++
			}

			for _, guest := range rsvp.Guests {
				if guest.Attending {
					data.AttendingGuests++
				}
			}
		}

		data.Rows = append(
			data.Rows,
			row,
		)
	}

	render(
		w,
		"admin.html",
		data,
	)
}

// --------------------------------------------------
// Helpers
// --------------------------------------------------

func isAdmin(r *http.Request) bool {
	cookie, err := r.Cookie("wedding_admin")

	return err == nil &&
		cookie.Value == "authenticated"
}

func getInvitation(
	code string,
) (Invitation, bool) {
	if repo != nil {
		return repo.FindInvitation(code)
	}

	code = normalizeCode(code)

	for _, invitation := range invitations {
		if invitation.Code == code {
			return invitation, true
		}
	}

	return Invitation{}, false
}

func invitationGuestNames(
	invitation Invitation,
) []string {
	parts := strings.Split(invitation.CoupleName, "&")
	if len(parts) == 1 {
		parts = strings.Split(invitation.CoupleName, "and")
	}

	result := make([]string, 0, 2)
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name != "" {
			result = append(result, name)
		}
	}

	if len(result) == 0 {
		return []string{invitation.CoupleName}
	}

	return result
}

func normalizeCode(
	code string,
) string {
	return strings.ToUpper(
		strings.TrimSpace(code),
	)
}

func urlQuery(
	value string,
) string {
	return template.URLQueryEscaper(value)
}
