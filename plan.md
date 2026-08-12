# Physiq — Planning Doc

## Overview
A weight and body measurement tracking app.

## Goals
- Make it easy to log weight and body measurements daily
- Show progress over time with clear trends and charts
- Keep data personal and private

## Scope

### In scope (v1)
- User accounts / login (SPA requires auth)
- Log weight, and body measurements (waist, chest, arms, etc.)
    - there will be some default/system measurements but the user can add custom measurements
- Timeline / history view of entries
- Charts showing trends over time
- Simple goals and progress tracking

### Out of scope (v1)
- Meal / calorie logging
- Workout plans
- Social features
- Wearable / smart-scale integration

## Decisions
- **Platform:** Web, responsive — all features must work in a mobile browser
- **PWA:** Out of scope for v1 (Android and iOS both support PWAs, but not building one yet)
- **Data entry:** Manual only for v1
- **Accounts:** Login required for the SPA; no cross-device sync in v1
- **Auth:** Session cookies (httpOnly) with Postgres-backed sessions; auth handled as pluggable middleware so OAuth/API keys can be added later for a public API
- **Units:** Per-user unit preference (weight: lb/kg, length: in/cm)
- **Security:** Password hashing (argon2/bcrypt) from day one; login rate limiting
- **Logging:** Structured logging in the backend
- **HTTPS:** Let's Encrypt for TLS in prod
- **CI/CD:** GitHub Actions (repo on GitHub) — lint, test, build
- **Testing:** Backend unit tests required in v1; expand coverage later

## Tech Stack
- **Frontend:** React SPA (interactive app: logging, history, charts)
- **Backend:** Go + Gin (JSON API for the SPA)
- **SSR pages:** Go + templ (landing / marketing pages)
- **Database:** PostgreSQL
- **Containerization:** Docker for dev and prod

## Future Imports/Integrations
- Smart scales (Withings, Garmin, Fitbit, Apple Health / Google Fit)
- Bulk CSV import/export
- Health app sync (HealthKit, Health Connect)
- Public API (OAuth / API keys — auth middleware already pluggable)

## Future Considerations
- **Testing expansion** — React component tests, e2e (e.g., Playwright)
- **Backups** — automated Postgres backups closer to launch
- **Observability** — metrics/monitoring

## Development & Docker

### Repo structure (monorepo)
```
physiq/
├── backend/           # Go + Gin + templ
│   ├── cmd/           # entrypoints
│   ├── internal/      # app code (handlers, middleware, db)
│   ├── templates/     # templ files
│   └── migrations/    # SQL migrations
├── frontend/          # React SPA
│   └── src/
├── docker-compose.yml # dev stack
├── compose.prod.yml   # prod stack
└── plan.md
```

### Dev
- `docker compose` stack with hot reload: Go API, Vite dev server (React), Postgres

```bash
docker compose up -d --build   # start dev stack (db + backend + frontend)
docker compose logs -f backend # follow backend (air hot reload) logs
docker compose down            # stop stack (keeps DB volume)
docker compose down -v         # stop + wipe DB volume
```

- App: http://localhost:5173 (React SPA, proxies /api → backend)
- API/landing: http://localhost:8080
- Backend regenerates templ + rebuilds automatically via air

### Prod
- Multi-stage Docker builds → minimal runtime images
- Go serves the built React SPA (static assets) + API; Postgres as a service
- Existing host nginx (already serving a Laravel app on the same server) acts as reverse proxy for physiq
- TLS: no Let's Encrypt setup exists yet — add certbot on the host for physiq's domain/subdomain
- Single-entry deployment (compose or similar)

```bash
cp .env.example .env          # set POSTGRES_USER/PASSWORD/DB
docker compose -f compose.prod.yml up -d --build   # runs migrate, then app
docker compose -f compose.prod.yml logs -f app
docker compose -f compose.prod.yml down
```

## Rough Milestones
1. **Setup** — repo structure, docker-compose dev stack, Postgres migrations
2. **Data model** — schema for users, entries, measurements (system + custom), goals
3. **Auth** — registration, login, sessions/tokens
4. **Landing pages** — Go + templ SSR marketing pages
5. **API** — CRUD endpoints for entries/measurements/goals
6. **React SPA** — login, logging entries, history view
7. **Visualization** — trend charts
8. **Goals** — set and track targets
9. **Polish** — testing, prod Docker build, deployment
