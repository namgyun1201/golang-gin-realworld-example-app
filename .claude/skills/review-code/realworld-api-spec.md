# RealWorld API Spec Reference

HTTP status codes, response formats, and endpoint contracts for review validation.

## Response Wrapper Format

All successful responses wrap data in a named key:

| Endpoint Type | Wrapper Key | Example |
|--------------|-------------|---------|
| Single user | `"user"` | `{"user": {"username": "...", "email": "...", "token": "...", "bio": "...", "image": "..."}}` |
| Profile | `"profile"` | `{"profile": {"username": "...", "bio": "...", "image": "...", "following": false}}` |
| Single article | `"article"` | `{"article": {"slug": "...", "title": "...", ...}}` |
| Article list | `"articles"` + `"articlesCount"` | `{"articles": [...], "articlesCount": 2}` |
| Single comment | `"comment"` | `{"comment": {"id": 1, "body": "...", ...}}` |
| Comment list | `"comments"` | `{"comments": [...]}` |
| Tags | `"tags"` | `{"tags": ["tag1", "tag2"]}` |

## Error Response Format

All errors use this structure:
```json
{
  "errors": {
    "body": ["error message 1", "error message 2"]
  }
}
```

Or with field-specific keys:
```json
{
  "errors": {
    "email": "{email: required}",
    "password": "{min: 8}"
  }
}
```

## HTTP Status Codes by Endpoint

### Users

| Endpoint | Success | Auth Failure | Validation | Not Found |
|----------|---------|-------------|------------|-----------|
| `POST /api/users` (register) | **201** Created | — | 422 | — |
| `POST /api/users/login` | 200 | **401** (NOT 403) | 422 | — |
| `GET /api/user` | 200 | 401 | — | — |
| `PUT /api/user` | 200 | 401 | 422 | — |

### Profiles

| Endpoint | Success | Auth Failure | Not Found |
|----------|---------|-------------|-----------|
| `GET /api/profiles/:username` | 200 | — (optional auth) | 404 |
| `POST /api/profiles/:username/follow` | 200 | 401 | 404 |
| `DELETE /api/profiles/:username/follow` | 200 | 401 | 404 |

### Articles

| Endpoint | Success | Auth Failure | Forbidden | Not Found |
|----------|---------|-------------|-----------|-----------|
| `GET /api/articles` | 200 | — (optional auth) | — | — |
| `GET /api/articles/feed` | 200 | 401 | — | — |
| `GET /api/articles/:slug` | 200 | — (optional auth) | — | 404 |
| `POST /api/articles` | **201** Created | 401 | — | — |
| `PUT /api/articles/:slug` | 200 | 401 | **403** | 404 |
| `DELETE /api/articles/:slug` | 200 (empty body) | 401 | **403** | 404 |

### Comments

| Endpoint | Success | Auth Failure | Forbidden | Not Found |
|----------|---------|-------------|-----------|-----------|
| `POST /api/articles/:slug/comments` | **201** Created | 401 | — | 404 |
| `GET /api/articles/:slug/comments` | 200 | — (optional auth) | — | 404 |
| `DELETE /api/articles/:slug/comments/:id` | 200 (empty body) | 401 | **403** | 404 |

### Favorites & Tags

| Endpoint | Success | Auth Failure | Not Found |
|----------|---------|-------------|-----------|
| `POST /api/articles/:slug/favorite` | 200 | 401 | 404 |
| `DELETE /api/articles/:slug/favorite` | 200 | 401 | 404 |
| `GET /api/tags` | 200 | — | — |

## Key Rules

1. **Login returns 401, NOT 403** — failed login is "Unauthorized", not "Forbidden"
2. **Register returns 201** — resource creation, not just success
3. **Delete endpoints return empty body** — no JSON response on successful delete
4. **Article/Comment creation returns 201** — consistent resource creation status
5. **Authorization failures (not owner) return 403** — distinct from authentication failure (401)
6. **Optional auth routes** — profiles and article list work without auth but return extra data with auth (e.g., `following` field)
7. **Pagination** — `GET /api/articles` uses `limit` (default 20) and `offset` (default 0) query params
8. **Filtering** — `GET /api/articles` supports `tag`, `author`, `favorited` query params
