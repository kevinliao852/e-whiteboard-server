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
- `HOST_AllOW_ORIGINS` as a comma-separated list of allowed frontend origins
- `DATABASE_PATH` if you do not want the default SQLite file path

Use a long random value for `SESSION_SECRET`.

Example:

```env
HOST_AllOW_ORIGINS=http://localhost:3000,http://127.0.0.1:3000,https://yourdomain.com
```

For remote frontend origins over HTTPS, the session cookie is issued as `SameSite=None; Secure` so the browser will send it on cross-origin API requests. For localhost development, the cookie stays `SameSite=Lax; Secure=false`.

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
