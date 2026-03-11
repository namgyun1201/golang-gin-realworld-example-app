# RealWorld API Endpoints Reference

Quick reference for E2E test validation.

## Endpoint Map

### Auth (no token required)
```
POST /api/users          → 201 { "user": {...} }
POST /api/users/login    → 200 { "user": {...} }
```

### User (token required)
```
GET  /api/user           → 200 { "user": {...} }
PUT  /api/user           → 200 { "user": {...} }
```

### Profiles (token optional for GET)
```
GET    /api/profiles/:username        → 200 { "profile": {...} }
POST   /api/profiles/:username/follow → 200 { "profile": {...} }
DELETE /api/profiles/:username/follow → 200 { "profile": {...} }
```

### Articles (token optional for GET, required for mutations)
```
GET    /api/articles            → 200 { "articles": [...], "articlesCount": N }
GET    /api/articles/feed       → 200 { "articles": [...], "articlesCount": N }
GET    /api/articles/:slug      → 200 { "article": {...} }
POST   /api/articles            → 201 { "article": {...} }
PUT    /api/articles/:slug      → 200 { "article": {...} }
DELETE /api/articles/:slug      → 200 (empty body)
```

### Comments (token optional for GET, required for mutations)
```
POST   /api/articles/:slug/comments      → 201 { "comment": {...} }
GET    /api/articles/:slug/comments      → 200 { "comments": [...] }
DELETE /api/articles/:slug/comments/:id  → 200 (empty body)
```

### Favorites (token required)
```
POST   /api/articles/:slug/favorite  → 200 { "article": {...} }
DELETE /api/articles/:slug/favorite  → 200 { "article": {...} }
```

### Tags (public)
```
GET /api/tags → 200 { "tags": [...] }
```

## Query Parameters

`GET /api/articles` supports:
- `tag=<string>` — filter by tag
- `author=<string>` — filter by author username
- `favorited=<string>` — filter by user who favorited
- `limit=<int>` — default 20
- `offset=<int>` — default 0

## Common Test Scenarios

1. **Register → Login → CRUD cycle**: create user, login, create article, update, delete
2. **Auth boundary**: access protected route without token → 401
3. **Authorization boundary**: update another user's article → 403
4. **Validation**: submit invalid data → 422 with field errors
5. **Pagination**: create multiple articles, verify limit/offset behavior
6. **Feed**: follow user, verify feed returns their articles
7. **Favorites**: favorite/unfavorite, verify count and `favorited` flag
