# Wedding RSVP — Go + Tailwind CSS

A small sample wedding website with:

- A pretty guest-facing landing page
- A unique invitation code for each invited couple
- RSVP form with attendance, guest count, food preference, allergies, and message
- Admin login and RSVP overview
- Responsive styling with Tailwind CSS

## Run the demo

Requirements:

- Go 1.22+
- Node.js + npm if you want to build Tailwind locally

### 1. Start the Go app

```bash
go run .
```

Open:

http://localhost:8080

### 2. Try the sample invitation codes

- `ALICE-BOB-7K2P`
- `CARLA-DAN-9M4Q`
- `EMMA-FELIX-3R8T`

You can also go directly to:

```text
http://localhost:8080/?code=ALICE-BOB-7K2P
```

## Language subdomains

The guest-facing pages select Czech, German, or English from the first hostname label:

- `cs.example.com` or `cz.example.com`: Czech
- `de.example.com`: German
- `en.example.com`: English

Point all three DNS records at the same deployment, configure a TLS certificate that covers all three hostnames, and route them to one Go app process. The reverse proxy must preserve the original `Host` header so the app can select the language. For example, a Caddy site can use:

```caddy
en.example.com, cs.example.com, de.example.com {
	reverse_proxy localhost:8080
}
```

Caddy can obtain HTTPS certificates automatically once DNS points to the server and ports 80/443 are reachable. With Nginx, forward the host using `proxy_set_header Host $host;`.

To test locally without DNS, send a host header:

```bash
curl -H 'Host: cs.example.com' http://localhost:8080/
```

Translations are in `i18n.go`; an unrecognized hostname defaults to English.

### 3. Admin

Open:

http://localhost:8080/admin

For this demo, the password is:

```text
change-me
```

For a real deployment, set:

```bash
export ADMIN_PASSWORD="a-long-random-password"
```

## Tailwind

The HTML currently loads Tailwind through the CDN so the demo works immediately.

If you want a locally compiled stylesheet instead:

```bash
npm install
npm run build:css
```

Then replace the Tailwind CDN `<script>` in `templates/base.html` with:

```html
<link rel="stylesheet" href="/static/css/app.css">
```

## Architecture and extensibility

The project is intentionally split so the public HTTP layer, domain model, storage concern, and template rendering are separated. That makes it easier to evolve the app without mixing request handling and persistence logic.

### Key files

- `main.go`: application startup and route registration.
- `models.go`: core invitation, RSVP, and page data structures.
- `store.go`: in-memory persistence plus a repository-style interface designed to be replaced by a database-backed implementation.
- `handlers.go`: request validation, form handling, and admin logic.
- `template.go`: template setup and rendering pipeline.
- `docs/architecture.md`: deeper technical notes for the future architecture.

This is a good foundation for moving to SQLite or PostgreSQL later, because the HTTP layer already depends on a storage contract rather than a concrete in-memory map.

## What to change for production

This is intentionally a sample, not a production-ready RSVP system.

Recommended next steps:

1. Replace the in-memory store with SQLite/PostgreSQL and migrate the current RSVP and invitation data model into tables.
2. Generate cryptographically random, non-guessable invitation tokens instead of manually assigned codes.
3. Add CSRF protection to forms.
4. Use proper admin authentication and secure session storage.
5. Add rate limiting to invitation-code and admin login endpoints.
6. Validate and constrain all submitted fields.
7. Add an export to CSV/Excel from the admin page.
8. Add email confirmation after an RSVP.
9. Put the site behind HTTPS.
10. Add a proper deployment configuration.

## Suggested future project structure

```text
wedding-rsvp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── handlers.go
│   │   └── template.go
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
├── docs/
│   └── architecture.md
├── README.md
└── go.mod
```

## Good next improvements

- Add a service layer for RSVP validation, guest counting, and summary reports.
- Add test coverage for edge cases like duplicate invitation codes, invalid codes, and admin authentication.
- Add database migrations and seed data.
- Add a nicer admin experience with filters, search, and CSV export.
- Add semantic logging and monitoring for production failures.
