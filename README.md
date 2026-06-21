# e-whiteboard

Backend server for the `e-whiteboard` app.

It provides:
- Google-based login
- whiteboard CRUD APIs
- websocket endpoints for collaborative drawing
- SQLite persistence for whiteboards and canvas data

![Image](https://i.imgur.com/XV8RHdK.gif)

## Setup

Copy the example environment file:

```bash
cp .env{.example,}
```

Set these values in `.env` before starting the server:
- `GOOGLE_CLIENT_ID`
- `SESSION_SECRET`
- `DATABASE_PATH` if you do not want the default SQLite file path

Use a long random value for `SESSION_SECRET`.

## Run

```bash
go run ./cmd
```

## Repos

The project is split across separate repositories:

| Name                         | Repo Address                                               |
|:-----------------------------|:-----------------------------------------------------------|
| e-whiteboard-server          | https://github.com/kevinliao852/e-whiteboard-server        |
| e-whiteboard-client          | https://github.com/kevinliao852/e-whiteboard-client        |
