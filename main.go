package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Invitation struct {
	Code       string
	CoupleName string
	Email      string
}

type Guest struct {
	Name      string
	Attending bool
	Allergies string
}

type RSVP struct {
	Code          string
	CoupleName    string
	Attending     bool
	GuestCount    int
	Accommodation string
	Allergies     string
	Message       string
	SubmittedAt   time.Time
	Guests        []Guest
}

type PageData struct {
	Invitation *Invitation
	RSVP       *RSVP
	HasRSVP    bool
	Error      string
}

type AdminData struct {
	TotalInvitations int
	Responses        int
	AttendingCouples int
	AttendingGuests  int
	Rows             []AdminRow
}

type AdminRow struct {
	Invitation Invitation
	RSVP       *RSVP
}

var (
	invitations = []Invitation{
		{
			Code:       "ALICE-BOB-7K2P",
			CoupleName: "Alice & Bob",
			Email:      "alice@example.com",
		},
		{
			Code:       "CARLA-DAN-9M4Q",
			CoupleName: "Carla & Dan",
			Email:      "carla@example.com",
		},
		{
			Code:       "EMMA-FELIX-3R8T",
			CoupleName: "Emma & Felix",
			Email:      "emma@example.com",
		},
	}

	rsvpStore = map[string]RSVP{}
	storeMu   sync.RWMutex

	tmpl = template.Must(
		template.New("base").Funcs(template.FuncMap{
			"invitationGuestNames": invitationGuestNames,
		}).ParseGlob("templates/*.html"),
	)
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", landingHandler)
	mux.HandleFunc("/rsvp", rsvpHandler)
	mux.HandleFunc("/admin/login", adminLoginHandler)
	mux.HandleFunc("/admin", adminHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf(
		"Wedding RSVP running on http://localhost:%s",
		port,
	)

	log.Fatal(
		http.ListenAndServe(":"+port, mux),
	)
}

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
		render(
			w,
			"landing.html",
			PageData{},
		)
		return
	}

	invitation, ok := getInvitation(code)

	if !ok {
		render(
			w,
			"landing.html",
			PageData{
				Error: "That invitation code was not found. Please check the code and try again.",
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

	storeMu.RLock()

	rsvp, hasRSVP := rsvpStore[invitation.Code]

	storeMu.RUnlock()

	var rsvpPtr *RSVP

	if hasRSVP {
		copy := rsvp
		rsvpPtr = &copy
	}

	render(
		w,
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
			render(
				w,
				"rsvp.html",
				PageData{
					Invitation: &invitation,
					Error:      "Please keep each guest's allergy note reasonably short.",
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
		render(
			w,
			"rsvp.html",
			PageData{
				Invitation: &invitation,
				Error:      "Please tell us how you would like to handle accommodation.",
			},
		)
		return
	}

	message := strings.TrimSpace(
		r.FormValue("message"),
	)

	if len(message) > 1000 {
		render(
			w,
			"rsvp.html",
			PageData{
				Invitation: &invitation,
				Error:      "Please keep your message reasonably short.",
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

	storeMu.Lock()

	rsvpStore[invitation.Code] = rsvp

	storeMu.Unlock()

	render(
		w,
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

	storeMu.RLock()
	defer storeMu.RUnlock()

	data := AdminData{
		TotalInvitations: len(invitations),
		Responses:        len(rsvpStore),
	}

	for _, invitation := range invitations {

		row := AdminRow{
			Invitation: invitation,
		}

		if rsvp, ok := rsvpStore[invitation.Code]; ok {

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

// Useful later when you create invitations dynamically.
func newInvitationCode() string {

	b := make([]byte, 6)

	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf(
		"%s-%s",
		strings.ToUpper(
			hex.EncodeToString(b[:3]),
		),
		strings.ToUpper(
			hex.EncodeToString(b[3:]),
		),
	)
}
