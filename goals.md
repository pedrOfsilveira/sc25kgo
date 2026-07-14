# SC25K Go Goals

## Product goal

Build a lightweight Couch to 5K app that is fast to start, simple to use, and easy
to maintain.

The first release should do one thing well: guide a runner through the complete
9-week, 27-session C25K program and remember their progress.

## Principles

- Ship one small Go binary with embedded static assets.
- Use SQLite for local persistence.
- Prefer the Go standard library and small dependencies.
- Use plain HTML, CSS, and minimal JavaScript; do not add a SPA framework.
- Keep the app usable on a phone and avoid unnecessary network requests.
- Add features only when they improve the core run experience.

## Current state

- [x] Go HTTP server using Chi
- [x] SQLite schema for users, stages, intervals, and completions
- [x] Stage list and stage detail queries
- [x] Basic completion and user queries
- [x] Three sample stages with interval data

These pieces are prototypes, not finished API features. Completion currently uses a
hard-coded user, photo, and XP value; XP is not added to the user; completed stages
are not user-scoped; and the user detail route is misspelled.

## Next: make the core reliable

- [x] Fix `/api/usser/{id}` to `/api/users/{id}`.
- [ ] Return `400` for invalid input, `404` for missing records, and consistent JSON errors for server failures.
- [ ] Add JSON field names to API models and return empty arrays as `[]`, not `null`.
- [ ] Remove hard-coded completion values; decode and validate request data.
- [ ] Make stage completion transactional and calculate rewards on the server.
- [ ] Scope completion history and progress to the active user.
- [ ] Enable SQLite foreign keys and add uniqueness/order constraints.
- [ ] Replace the checked-in development database with repeatable schema migrations and idempotent seed data.
- [ ] Add graceful shutdown, server timeouts, and environment-based address/database configuration.

## MVP: complete workout loop

### Training plan

- [ ] Seed all 9 weeks and 27 sessions of the C25K plan.
- [ ] Store every warm-up, walk, run, and cool-down interval in execution order.
- [ ] Validate that seeded workouts have valid interval types and durations.
- [ ] Show session duration and interval summary before a run starts.

### Run experience

- [ ] Add a responsive home screen showing weeks, sessions, and progress.
- [ ] Add a workout screen with the current interval, countdown, next interval, and
      total progress.
- [ ] Support start, pause, resume, skip warm-up, and cancel.
- [ ] Keep the timer correct when the browser is backgrounded by deriving remaining
      time from timestamps rather than counting ticks.
- [ ] Add optional sound and vibration cues with a mute setting.
- [ ] Confirm before abandoning an active workout.

### Progress

- [ ] Use one local profile for the MVP; no sign-up or authentication required.
- [ ] Record completed and canceled runs with elapsed time and completion percentage.
- [ ] Award full XP for a first completion, 50% XP for repeats, and proportional XP
      for canceled runs.
- [ ] Update XP and insert the run record in one database transaction.
- [ ] Show the next recommended session and a simple run history.
- [ ] Allow completed sessions to be repeated without changing program progress.

## API target

- [x] `GET /api/stages` — list sessions
- [x] `GET /api/stages/{id}` — get a session and its intervals
- [ ] `GET /api/progress` — get profile, next session, XP, and completion summary
- [ ] `GET /api/runs` — list run history
- [ ] `POST /api/stages/{id}/complete` — record a validated run result
- [ ] `PATCH /api/profile` — update the local display name and preferences
- [ ] `GET /healthz` — report process and database health

The current `/api/users` endpoints are not part of the single-profile MVP. Remove
them unless multi-user support becomes an explicit requirement.

## Quality and release

- [ ] Add database tests using a temporary SQLite database.
- [ ] Add handler tests with `httptest` for success, validation, and not-found cases.
- [ ] Test reward rules, repeat runs, canceled runs, and transaction rollback.
- [ ] Run `go test ./...`, `go vet ./...`, and a production build in CI.
- [ ] Embed versioned migrations, seed data, and web assets with `go:embed`.
- [ ] Ensure the app starts with an empty database and requires no manual setup.
- [ ] Keep the release artifact to one executable plus one writable database file.
- [ ] Document local development, configuration, backup, and upgrade steps.

## Definition of MVP done

- A new install automatically contains all 27 sessions.
- A runner can complete an entire guided workout from a phone.
- Closing or backgrounding the page does not corrupt the timer or progress.
- Completion, repeat, cancellation, and XP rules are deterministic and tested.
- Progress survives a restart and identifies the correct next session.
- The app builds and runs as a single small Go service with embedded web assets.

## Deliberately deferred

These features from the original SC25K are not required for the lightweight first
release:

- Accounts, password recovery, and cloud sync
- Global rankings and social features
- Shop, offers, and XP purchases
- Avatars and photo uploads
- Badges, confetti, and shareable cards
- Push notifications and analytics
- Advanced offline/background sync
- Multiple languages

Reconsider deferred features only after the core 27-session journey is reliable and
people are using it.
