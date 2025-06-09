package simplemailtp

import (
	"testing"

	"mtg-report/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	mockAuth := mocks.NewSMTPAuthMock()
	mockTimer := mocks.NewTimerMock()

	from := "test@example.com"
	to := "user@example.com"
	address := "smtp.example.com:587"

	emailService := New(mockAuth, mockTimer, from, to, address)

	assert.NotNil(t, emailService)
	assert.Equal(t, mockAuth, emailService.auth)
	assert.Equal(t, mockTimer, emailService.timer)
	assert.Equal(t, from, emailService.from)
	assert.Equal(t, to, emailService.to)
	assert.Equal(t, address, emailService.adress)
}

func TestSendEmail_MessageConstruction(t *testing.T) {
	// This test focuses on verifying the email service constructor
	// and basic functionality without actually sending emails
	// since smtp.SendMail is a global function that's hard to mock

	mockAuth := mocks.NewSMTPAuthMock()
	mockTimer := mocks.NewTimerMock()

	from := "test@example.com"
	to := "user@example.com"
	address := "invalid-address" // Using invalid address to prevent actual sending

	emailService := New(mockAuth, mockTimer, from, to, address)

	cardsTable := "<table><tr><td>Test Card</td></tr></table>"
	cardsPrice := "Total price increased from $100 to $150"

	mockTimer.On("Now").Return("2023-12-01 10:00:00")

	// This will fail because of invalid address, but we can verify
	// the service was constructed properly and timer was called
	err := emailService.SendEmail(cardsTable, cardsPrice)

	// We expect an error because of invalid address
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")

	// Verify timer was called
	mockTimer.AssertExpectations(t)
}
