package tgbot

import (
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUserFromUpdate(t *testing.T) {
	tests := []struct {
		name     string
		update   *models.Update
		expected *models.User
	}{
		{
			name: "Message with From",
			update: &models.Update{
				Message: &models.Message{
					From: &models.User{ID: 123, Username: "testuser"},
				},
			},
			expected: &models.User{ID: 123, Username: "testuser"},
		},
		{
			name: "EditedMessage with From",
			update: &models.Update{
				EditedMessage: &models.Message{
					From: &models.User{ID: 456, Username: "editeduser"},
				},
			},
			expected: &models.User{ID: 456, Username: "editeduser"},
		},
		{
			name:     "No user information",
			update:   &models.Update{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUserFromUpdate(tt.update)
			assert.Equal(t, tt.expected, result)
		})
	}
}
