package rd_test

import (
	"testing"
	"time"

	"github.com/nenadstojanovikj/rd"
	"github.com/stretchr/testify/assert"
)

func TestToken_ExpiredAuthorization(t *testing.T) {
	token := rd.Token{
		AccessToken:  "ACCESS",
		RefreshToken: "REFRESH",
		ObtainedAt:   time.Now().Add(-3600 * time.Second),
	}

	assert.False(t, token.IsValid())
}
