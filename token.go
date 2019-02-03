package rd

import (
	"time"
)

type TokenRefresher interface {
	RefreshAccessToken(token Token) (t Token, err error)
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ObtainedAt   time.Time
}

// IsValid checks if the current token is not expired and valid
func (t *Token) IsValid() bool {
	// We expire the token 10 seconds before, so we don't send a request and risk to have it failing during transport
	tokenExpiry := t.ObtainedAt.Add(time.Second * time.Duration(t.ExpiresIn-10))
	return time.Now().Before(tokenExpiry)
}
