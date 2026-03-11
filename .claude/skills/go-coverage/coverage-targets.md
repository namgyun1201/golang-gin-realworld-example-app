# Coverage Targets & Strategy

## Package Targets

| Package | Target | Priority Areas |
|---------|--------|---------------|
| `users` | 99%+ | Auth middleware, password hashing, JWT flow, profile follow/unfollow |
| `articles` | 93%+ | CRUD handlers, authorization checks, pagination, favorites, comments |
| `common` | 85%+ | DB init, JWT generation/validation, error formatting, validators |
| **Overall** | **90%+** | — |

## Test Infrastructure

- Test DB setup: `common.TestDBInit()` → creates isolated SQLite test DB
- Test DB teardown: `common.TestDBFree()` → closes connection, deletes file
- Test files: `*_test.go` co-located with source files
- Assertions: `stretchr/testify` (`assert`, `require`)

## High-Value Test Targets

### users package
- `AuthMiddleware` — valid token, invalid token, expired token, missing token, optional auth
- `UsersRegistration` — success, duplicate email, validation failures
- `UsersLogin` — success, wrong email, wrong password, validation failures
- `UserUpdate` — partial update, password change, validation failures
- `ProfileRetrieve` — with auth (following status), without auth
- `Follow/Unfollow` — success, self-follow prevention, not found

### articles package
- `ArticlesCRUD` — create, read, update, delete with authorization checks
- `ArticlesList` — pagination (limit/offset), filtering (tag/author/favorited)
- `ArticlesFeed` — only followed authors' articles
- `Comments` — create, list, delete with authorization
- `Favorites` — favorite/unfavorite, count accuracy, `favorited` flag per user

### common package
- `GenToken/ParseToken` — round-trip, expiry, invalid signatures
- `Bind` — different content types, validation errors
- `NewValidatorError` — field name extraction, tag formatting
- `RandString/RandInt` — uniqueness, length, character set
