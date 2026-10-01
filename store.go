package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// RSVPRepository defines the persistence contract used by the application.
// This keeps the HTTP handlers independent from the concrete storage layer.
type RSVPRepository interface {
	ListInvitations() []Invitation
	FindInvitation(code string) (Invitation, bool)
	SaveRSVP(rsvp RSVP) error
	GetRSVP(code string) (RSVP, bool)
}

// MemoryRepository is the in-memory fallback used for demos and tests.
type MemoryRepository struct {
	mu          sync.RWMutex
	invitations []Invitation
	responses   map[string]RSVP
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		invitations: []Invitation{
			{Code: "ALICE-BOB-7K2P", CoupleName: "Alice & Bob", Email: "alice@example.com"},
			{Code: "CARLA-DAN-9M4Q", CoupleName: "Carla & Dan", Email: "carla@example.com"},
			{Code: "EMMA-FELIX-3R8T", CoupleName: "Emma & Felix", Email: "emma@example.com"},
		},
		responses: map[string]RSVP{},
	}
}

func (r *MemoryRepository) ListInvitations() []Invitation {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Invitation, len(r.invitations))
	copy(out, r.invitations)
	return out
}

func (r *MemoryRepository) FindInvitation(code string) (Invitation, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, invitation := range r.invitations {
		if invitation.Code == code {
			return invitation, true
		}
	}

	return Invitation{}, false
}

func (r *MemoryRepository) SaveRSVP(rsvp RSVP) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.responses[rsvp.Code] = rsvp
	return nil
}

func (r *MemoryRepository) GetRSVP(code string) (RSVP, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rsvp, ok := r.responses[code]
	return rsvp, ok
}

// SQLiteRepository persists invitations and RSVPs to a SQLite database.
type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dataSourceName string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	repo := &SQLiteRepository{db: db}
	if err := repo.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return repo, nil
}

func (r *SQLiteRepository) migrate() error {
	if _, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS invitations (
			code TEXT PRIMARY KEY,
			couple_name TEXT NOT NULL,
			email TEXT NOT NULL
		);
	`); err != nil {
		return fmt.Errorf("create invitations table: %w", err)
	}

	if _, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS rsvps (
			code TEXT PRIMARY KEY,
			couple_name TEXT NOT NULL,
			attending INTEGER NOT NULL,
			guest_count INTEGER NOT NULL,
			accommodation TEXT NOT NULL,
			allergies TEXT NOT NULL,
			message TEXT NOT NULL,
			submitted_at TEXT NOT NULL,
			guests TEXT NOT NULL
		);
	`); err != nil {
		return fmt.Errorf("create rsvps table: %w", err)
	}

	seed := []Invitation{
		{Code: "ALICE-BOB-7K2P", CoupleName: "Alice & Bob", Email: "alice@example.com"},
		{Code: "CARLA-DAN-9M4Q", CoupleName: "Carla & Dan", Email: "carla@example.com"},
		{Code: "EMMA-FELIX-3R8T", CoupleName: "Emma & Felix", Email: "emma@example.com"},
	}

	for _, invitation := range seed {
		if _, err := r.db.Exec(
			`INSERT OR IGNORE INTO invitations (code, couple_name, email) VALUES (?, ?, ?)`,
			invitation.Code,
			invitation.CoupleName,
			invitation.Email,
		); err != nil {
			return fmt.Errorf("seed invitation %s: %w", invitation.Code, err)
		}
	}

	return nil
}

func (r *SQLiteRepository) ListInvitations() []Invitation {
	rows, err := r.db.Query(`SELECT code, couple_name, email FROM invitations ORDER BY code`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var invitations []Invitation
	for rows.Next() {
		var invitation Invitation
		if err := rows.Scan(&invitation.Code, &invitation.CoupleName, &invitation.Email); err != nil {
			continue
		}
		invitations = append(invitations, invitation)
	}

	return invitations
}

func (r *SQLiteRepository) FindInvitation(code string) (Invitation, bool) {
	var invitation Invitation

	err := r.db.QueryRow(
		`SELECT code, couple_name, email FROM invitations WHERE code = ?`,
		normalizeCode(code),
	).Scan(&invitation.Code, &invitation.CoupleName, &invitation.Email)
	if err != nil {
		return Invitation{}, false
	}

	return invitation, true
}

func (r *SQLiteRepository) SaveRSVP(rsvp RSVP) error {
	guestBytes, err := json.Marshal(rsvp.Guests)
	if err != nil {
		return fmt.Errorf("marshal guests: %w", err)
	}

	_, err = r.db.Exec(`
		INSERT INTO rsvps (
			code,
			couple_name,
			attending,
			guest_count,
			accommodation,
			allergies,
			message,
			submitted_at,
			guests
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(code) DO UPDATE SET
			couple_name = excluded.couple_name,
			attending = excluded.attending,
			guest_count = excluded.guest_count,
			accommodation = excluded.accommodation,
			allergies = excluded.allergies,
			message = excluded.message,
			submitted_at = excluded.submitted_at,
			guests = excluded.guests
	`,
		rsvp.Code,
		rsvp.CoupleName,
		boolToInt(rsvp.Attending),
		rsvp.GuestCount,
		rsvp.Accommodation,
		rsvp.Allergies,
		rsvp.Message,
		rsvp.SubmittedAt.UTC().Format(time.RFC3339Nano),
		string(guestBytes),
	)
	if err != nil {
		return fmt.Errorf("save rsvp: %w", err)
	}

	return nil
}

func (r *SQLiteRepository) GetRSVP(code string) (RSVP, bool) {
	var (
		rsvp        RSVP
		attending   int
		submittedAt string
		guestsJSON  string
	)

	err := r.db.QueryRow(`
		SELECT code, couple_name, attending, guest_count, accommodation, allergies, message, submitted_at, guests
		FROM rsvps WHERE code = ?
	`, normalizeCode(code)).Scan(
		&rsvp.Code,
		&rsvp.CoupleName,
		&attending,
		&rsvp.GuestCount,
		&rsvp.Accommodation,
		&rsvp.Allergies,
		&rsvp.Message,
		&submittedAt,
		&guestsJSON,
	)
	if err != nil {
		return RSVP{}, false
	}

	rsvp.Attending = attending == 1
	if parsed, err := time.Parse(time.RFC3339Nano, submittedAt); err == nil {
		rsvp.SubmittedAt = parsed
	}
	if err := json.Unmarshal([]byte(guestsJSON), &rsvp.Guests); err != nil {
		rsvp.Guests = nil
	}

	return rsvp, true
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

var (
	repo        = mustOpenRepository()
	invitations = repo.ListInvitations()
	rsvpStore   = map[string]RSVP{}
	storeMu     sync.RWMutex
)

func mustOpenRepository() RSVPRepository {
	dataSourceName := os.Getenv("DATABASE_URL")
	if dataSourceName == "" {
		dataSourceName = "wedding-rsvp.db"
	}

	repository, err := NewSQLiteRepository(dataSourceName)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize sqlite repository: %v", err))
	}

	return repository
}
