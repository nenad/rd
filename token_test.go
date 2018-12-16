package realdebrid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToken_ExpiredAuthorization(t *testing.T) {
	token := Token{
		AccessToken:  "ACCESS",
		RefreshToken: "REFRESH",
		ObtainedAt:   time.Now().Add(-3600 * time.Second),
	}

	assert.False(t, token.IsValid())
}
