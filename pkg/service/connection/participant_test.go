package connection

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewParticipant(t *testing.T) {

	t.Run("initializer test", func(t *testing.T) {
		p := NewParticipant()

		assert.NotEmpty(t, p.ParticipantID)
		assert.NotEmpty(t, p.PK)
		assert.True(t, strings.HasPrefix(p.PK, "participantID:"))
		assert.True(t, strings.HasPrefix(p.ParticipantID, "p-"))
		expiredAt := time.Unix(p.ExpiredAt, 0)
		limit := time.Since(expiredAt).Hours()
		assert.Less(t, limit, 24.0)
	})

	t.Run("first connect", func(t *testing.T) {
		p := NewParticipant()
		p.P256dh = "ut-curve"
		p.Auth = "ut-auth"
		p.PK = "ut:ut-first-connect"
		p.Save()
	})

	t.Run("retreive participant", func(t *testing.T) {
		testParticipantId := "ut-find-test"
		p := NewParticipant()
		p.ParticipantID = testParticipantId
		p.PK = pkPrefixParticipantId + testParticipantId
		p.Save()

		found, err := FindParticipant(testParticipantId)
		assert.Nil(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, testParticipantId, found.ParticipantID)
	})
}
