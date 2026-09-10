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

## What to change for production

This is intentionally a sample, not a production-ready RSVP system.

Recommended next steps:

1. Move invitations and RSVPs from in-memory maps into SQLite/PostgreSQL.
2. Generate cryptographically random, non-guessable invitation tokens instead of manually assigned codes.
3. Add CSRF protection to forms.
4. Use proper admin authentication and secure session storage.
5. Add rate limiting to invitation-code and admin login endpoints.
6. Validate and constrain all submitted fields.
7. Add an export to CSV/Excel from the admin page.
8. Add email confirmation after an RSVP.
9. Put the site behind HTTPS.
10. Add a proper deployment configuration.

## Suggested project structure

```text
wedding-rsvp/
├── main.go
├── go.mod
├── package.json
├── tailwind.config.js
├── templates/
│   ├── base.html
│   ├── landing.html
│   ├── rsvp.html
│   ├── success.html
│   ├── admin-login.html
│   └── admin.html
└── static/
    └── css/
        └── input.css
```
