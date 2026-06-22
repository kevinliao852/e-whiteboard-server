# Guest Login Spec

## Roles

- `guest`
- `user`

## Authentication Endpoints

- `POST /guest-login`
- `POST /login`
- `GET /v1/me`

## Behavior

- Guest users are authenticated through the same session-cookie mechanism as registered users.
- Guest users may join and use whiteboards.
- Guest users must not be allowed to create or delete whiteboards.

## Authorized Actions

### Guest

- `GET /v1/me`
- `GET /v1/whiteboards`
- `GET /v1/whiteboards/:id/points`
- `GET /v1/chat-messages?room-id=:id`
- `WS /ws/chat/:id`
- `WS /ws/drawing/:id`

### Registered User

- All guest actions
- `POST /v1/whiteboards`
- `DELETE /v1/whiteboards/:id`

## Disallowed Actions For Guests

- Create whiteboard
- Delete whiteboard

## Response Expectations

- `GET /v1/me` must indicate the current role.
- Guest session responses must identify the caller as `guest`.
- Registered session responses must identify the caller as `user`.

## UI Requirements

- Guest users must not see enabled create or delete whiteboard actions.
- Guest users must be able to open existing whiteboards and collaborate normally.
- The frontend must treat the session cookie as the source of truth for auth state.

