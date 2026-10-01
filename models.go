package main

import "time"

// Invitation represents the invitee information tied to a unique RSVP code.
// In a production system this data would likely come from a database table
// such as invitations with fields like id, code, couple_name, email, created_at.
type Invitation struct {
	Code       string
	CoupleName string
	Email      string
}

// Guest represents one guest on an invitation.
type Guest struct {
	Name      string
	Attending bool
	Allergies string
}

// RSVP stores the full response for a single invitation code.
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

// PageData is used for rendering the public-facing pages.
type PageData struct {
	Invitation *Invitation
	RSVP       *RSVP
	HasRSVP    bool
	Error      string
}

// AdminData is used by the admin dashboard.
type AdminData struct {
	TotalInvitations int
	Responses        int
	AttendingCouples int
	AttendingGuests  int
	Rows             []AdminRow
}

// AdminRow groups a single invitation with its optional RSVP response.
type AdminRow struct {
	Invitation Invitation
	RSVP       *RSVP
}
