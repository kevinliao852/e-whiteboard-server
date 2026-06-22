# Folder Structure Guideline

This project follows a clean-architecture oriented layout.

## Core Idea

- `core` contains domain entities and port interfaces.
- `service` contains use-case logic and orchestration.
- `adapter` contains concrete implementations for the outside world.
- `route` contains application wiring and route registration only.

## Recommended Layout

```text
cmd/
  main.go

internal/
  core/
    - entities
    - service/repository ports

  service/
    - application use cases
    - orchestration

  adapter/
    db/
      - repository implementations
      - GORM / SQL / persistence code

    web/
      controllers/
        - HTTP handlers
        - WebSocket handlers
      middlewares/
        - HTTP middleware
      ws/
        - websocket connection helpers
        - websocket message types

    state/
      - in-memory room registry
      - live websocket room state
      - transient connection tracking

  route/
    - composition root
    - dependency injection
    - route registration

  database/
    - DB connection setup
```

## Package Responsibilities

### `core`

- Keep business concepts here.
- Define interfaces that outer layers implement.
- Do not import Gin, GORM, websocket libraries, or other infrastructure packages.

### `service`

- Put application rules here.
- Coordinate persistence, state, and delivery-layer input.
- Keep framework-specific code out of this package.

### `adapter/db`

- Put repository implementations here.
- Map core types to database records.
- Keep all GORM and SQL logic here.

### `adapter/web`

- Put HTTP and WebSocket delivery code here.
- Handle request parsing, response formatting, and websocket connection management.
- Keep business rules out of controllers and middleware.

### `adapter/state`

- Put in-memory runtime state here.
- Use it for live websocket rooms, participant tracking, and other transient process-local state.
- Do not use it for durable business data.

### `route`

- Keep this package focused on wiring.
- Instantiate services, repositories, registries, and controllers here.
- Do not put business logic here.

## Naming Rules

- Use `adapter/db` for repository implementations.
- Use `adapter/web` for HTTP and WebSocket adapters.
- Use `adapter/state` for in-memory runtime state.
- Use `route` or `bootstrap` for composition only.

## Practical Rules

- If data must survive restart, store it in `adapter/db`.
- If data only exists while the process is running, store it in `adapter/state`.
- If code talks to HTTP clients or websocket clients, keep it in `adapter/web`.
- If code decides what the app should do, keep it in `service`.
