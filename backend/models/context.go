package models

// contextKey is a private type to prevent key collisions across packages.
type contextKey string

// UserIDKey is the key used to store the authenticated user's ID in
// the request context after JWT verification.
const UserIDKey contextKey = "userID"