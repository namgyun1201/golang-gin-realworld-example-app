# Security Checklist

Patterns to verify during code review for this project.

## Authentication

- [ ] Protected routes use `AuthMiddleware(true)` (auto401=true)
- [ ] Optional auth routes use `AuthMiddleware(false)` (auto401=false)
- [ ] JWT token extracted from `Authorization: Token <jwt>` header
- [ ] JWT signing uses HMAC-SHA256 (`jwt.SigningMethodHS256`)
- [ ] Token expiry set (24 hours from creation)
- [ ] Token validation checks signing method before accepting

## Authorization

- [ ] Update/delete article checks `article.AuthorID == currentUserID`
- [ ] Delete comment checks `comment.AuthorID == currentUserID`
- [ ] User update only modifies the authenticated user's own data
- [ ] Authorization failure returns 403 Forbidden (not 401 or 404)

## Input Validation

- [ ] All request bodies validated via `common.Bind()` with validator tags
- [ ] Username: `required,min=4,max=255`
- [ ] Email: `required,email`
- [ ] Password: `required,min=8,max=255`
- [ ] Bio: `max=1024`
- [ ] Image: `omitempty,url`
- [ ] Article title: `required,min=4`
- [ ] Article description: `required,max=2048`
- [ ] Article body: `required,max=2048`
- [ ] Comment body: `required,max=2048`
- [ ] Using `required` tag (NOT deprecated `exists`)

## Password Security

- [ ] Passwords hashed with `bcrypt.GenerateFromPassword` at `bcrypt.DefaultCost`
- [ ] Password verification uses `bcrypt.CompareHashAndPassword`
- [ ] Empty passwords rejected before hashing
- [ ] Password field never included in serialized responses
- [ ] Update uses `RandomPassword` sentinel to skip unnecessary re-hashing

## Database Security

- [ ] All queries use GORM (parameterized, no raw SQL concatenation)
- [ ] No `db.Raw()` or `db.Exec()` with string concatenation
- [ ] Transaction used for multi-step operations

## Error Handling

- [ ] Database errors don't leak internal details to client
- [ ] Authentication errors use generic message ("Not Registered email or invalid password")
- [ ] Validation errors use structured format (`errors` object)
- [ ] No panic in request handlers
