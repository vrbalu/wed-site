# Architecture notes

This project is intentionally small, but the code is structured to be easy to move toward a real application architecture.

## Current responsibilities

- `main.go`: bootstraps the HTTP server and runtime configuration.
- `models.go`: defines the domain entities used across the app.
- `store.go`: holds the in-memory storage and exposes a repository-shaped interface.
- `handlers.go`: contains HTTP request handling and validation logic.
- `template.go`: configures the HTML template engine and rendering pipeline.
- `templates/*.html`: presentation layer for the landing, RSVP, and admin pages.

## Why this is a good base for a database

The main stability point is the storage abstraction. The code already separates the core domain structs from the request layer and the in-memory implementation, which means a future repository can be swapped in without rewriting the handlers.

A likely migration path is:

1. Replace the in-memory map with SQLite or PostgreSQL.
2. Persist invitations and RSVP responses using repositories that satisfy the same interface.
3. Keep the HTTP handlers mostly untouched while changing only the persistence implementation.

## Recommended next steps

- Add a real `InvitationRepository` and `RSVPRepository` backed by SQL.
- Move validation into domain/service functions rather than keeping it in the HTTP layer.
- Add structured logging and request IDs.
- Introduce a config package for environment variables and secrets.
- Add CSRF tokens and secure session management for admin auth.
- Add email delivery or webhook hooks after an RSVP submission.

## Suggested future structure

```text
wedding-rsvp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── handlers.go
│   │   └── templates.go
│   ├── domain/
│   │   └── models.go
│   ├── repository/
│   │   ├── interface.go
│   │   ├── memory.go
│   │   └── postgres.go
│   └── service/
│       └── rsvp_service.go
├── templates/
├── static/
└── README.md
```

This layout keeps the app easier to test and makes the database layer easier to replace as the project grows.
