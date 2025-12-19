package auth

import "testing"

func TestAuthMiddleware(t *testing.T) {
	_ = AuthMiddleware("secret")
}
