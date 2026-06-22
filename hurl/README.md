# Hurl examples

Replace the placeholder values before running the requests.

## Files

- `user_login.hurl` - login with a Google ID token
- `get_user.hurl` - fetch a user by numeric ID
- `list_whiteboards.hurl` - list whiteboard IDs for a user
- `create_whiteboard.hurl` - create a new whiteboard
- `delete_whiteboard.hurl` - delete a whiteboard

## Notes

- The server defaults to `http://localhost:8080`
- `GET /v1/user/:id` may require an authenticated session depending on middleware config
- `DELETE /v1/whiteboards/:id` expects a JSON body with `whiteboard_id`
