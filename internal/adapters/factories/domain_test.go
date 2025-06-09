package factories

import (
	"mtg-report/internal/adapters/entities"
	"mtg-report/internal/core/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCardsInfoToCardsDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    []entities.MysqlCardInfo
		expected []domain.Cards
	}{
		{
			name:     "should return empty slice when input is empty",
			input:    []entities.MysqlCardInfo{},
			expected: []domain.Cards{},
		},
		{
			name: "should convert single card info with nil last price",
			input: []entities.MysqlCardInfo{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					LastPrice:       nil,
				},
			},
			expected: []domain.Cards{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					CardsDetails: domain.CardsDetails{
						LastPrice: 0,
					},
				},
			},
		},
		{
			name: "should convert single card info with last price",
			input: []entities.MysqlCardInfo{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					LastPrice:       floatPtr(15.50),
				},
			},
			expected: []domain.Cards{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					CardsDetails: domain.CardsDetails{
						LastPrice: 15.50,
					},
				},
			},
		},
		{
			name: "should convert multiple card infos",
			input: []entities.MysqlCardInfo{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					LastPrice:       floatPtr(15.50),
				},
				{
					ID:              2,
					Name:            "Counterspell",
					SetName:         "M21",
					CollectorNumber: "456",
					Foil:            false,
					LastPrice:       nil,
				},
				{
					ID:              3,
					Name:            "Dark Ritual",
					SetName:         "M21",
					CollectorNumber: "789",
					Foil:            true,
					LastPrice:       floatPtr(25.00),
				},
			},
			expected: []domain.Cards{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					CardsDetails: domain.CardsDetails{
						LastPrice: 15.50,
					},
				},
				{
					ID:              2,
					Name:            "Counterspell",
					SetName:         "M21",
					CollectorNumber: "456",
					Foil:            false,
					CardsDetails: domain.CardsDetails{
						LastPrice: 0,
					},
				},
				{
					ID:              3,
					Name:            "Dark Ritual",
					SetName:         "M21",
					CollectorNumber: "789",
					Foil:            true,
					CardsDetails: domain.CardsDetails{
						LastPrice: 25.00,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CardsInfoToCardsDomain(tt.input)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, len(tt.input), len(result))
		})
	}
}

func TestCardPriceHistoryToCardsDomain(t *testing.T) {
	fixedTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		input    []entities.MysqlCardPriceHistory
		expected []domain.Cards
	}{
		{
			name:     "should return empty slice when input is empty",
			input:    []entities.MysqlCardPriceHistory{},
			expected: []domain.Cards{},
		},
		{
			name: "should convert single card price history with nil last update",
			input: []entities.MysqlCardPriceHistory{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					LastPrice:       15.50,
					OldPrice:        12.00,
					PriceChange:     3.50,
					LastUpdate:      nil,
				},
			},
			expected: []domain.Cards{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					CardsDetails: domain.CardsDetails{
						LastPrice:   15.50,
						OldPrice:    12.00,
						PriceChange: 3.50,
						LastUpdate:  &time.Time{},
					},
				},
			},
		},
		{
			name: "should convert single card price history with last update",
			input: []entities.MysqlCardPriceHistory{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					LastPrice:       15.50,
					OldPrice:        12.00,
					PriceChange:     3.50,
					LastUpdate:      &fixedTime,
				},
			},
			expected: []domain.Cards{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					CardsDetails: domain.CardsDetails{
						LastPrice:   15.50,
						OldPrice:    12.00,
						PriceChange: 3.50,
						LastUpdate:  &fixedTime,
					},
				},
			},
		},
		{
			name: "should convert multiple card price histories",
			input: []entities.MysqlCardPriceHistory{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					LastPrice:       15.50,
					OldPrice:        12.00,
					PriceChange:     3.50,
					LastUpdate:      &fixedTime,
				},
				{
					ID:              2,
					Name:            "Counterspell",
					SetName:         "M21",
					CollectorNumber: "456",
					Foil:            false,
					LastPrice:       8.00,
					OldPrice:        10.00,
					PriceChange:     -2.00,
					LastUpdate:      nil,
				},
				{
					ID:              3,
					Name:            "Dark Ritual",
					SetName:         "M21",
					CollectorNumber: "789",
					Foil:            true,
					LastPrice:       25.00,
					OldPrice:        20.00,
					PriceChange:     5.00,
					LastUpdate:      &fixedTime,
				},
			},
			expected: []domain.Cards{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					CardsDetails: domain.CardsDetails{
						LastPrice:   15.50,
						OldPrice:    12.00,
						PriceChange: 3.50,
						LastUpdate:  &fixedTime,
					},
				},
				{
					ID:              2,
					Name:            "Counterspell",
					SetName:         "M21",
					CollectorNumber: "456",
					Foil:            false,
					CardsDetails: domain.CardsDetails{
						LastPrice:   8.00,
						OldPrice:    10.00,
						PriceChange: -2.00,
						LastUpdate:  &time.Time{},
					},
				},
				{
					ID:              3,
					Name:            "Dark Ritual",
					SetName:         "M21",
					CollectorNumber: "789",
					Foil:            true,
					CardsDetails: domain.CardsDetails{
						LastPrice:   25.00,
						OldPrice:    20.00,
						PriceChange: 5.00,
						LastUpdate:  &fixedTime,
					},
				},
			},
		},
		{
			name: "should handle zero values correctly",
			input: []entities.MysqlCardPriceHistory{
				{
					ID:              0,
					Name:            "",
					SetName:         "",
					CollectorNumber: "",
					Foil:            false,
					LastPrice:       0.0,
					OldPrice:        0.0,
					PriceChange:     0.0,
					LastUpdate:      nil,
				},
			},
			expected: []domain.Cards{
				{
					ID:              0,
					Name:            "",
					SetName:         "",
					CollectorNumber: "",
					Foil:            false,
					CardsDetails: domain.CardsDetails{
						LastPrice:   0.0,
						OldPrice:    0.0,
						PriceChange: 0.0,
						LastUpdate:  &time.Time{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CardPriceHistoryToCardsDomain(tt.input)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, len(tt.input), len(result))
		})
	}
}

func TestCardsInfoToCardsDomain_Capacity(t *testing.T) {
	// Test that the slice is created with the correct capacity
	input := make([]entities.MysqlCardInfo, 100)
	for i := 0; i < 100; i++ {
		input[i] = entities.MysqlCardInfo{
			ID:              int64(i + 1),
			Name:            "Card",
			SetName:         "SET",
			CollectorNumber: "123",
			Foil:            i%2 == 0,
			LastPrice:       floatPtr(float64(i + 1)),
		}
	}

	result := CardsInfoToCardsDomain(input)
	assert.Equal(t, 100, len(result))
	assert.Equal(t, 100, cap(result))
}

func TestCardPriceHistoryToCardsDomain_Capacity(t *testing.T) {
	// Test that the slice is created with the correct capacity
	fixedTime := time.Now()
	input := make([]entities.MysqlCardPriceHistory, 50)
	for i := 0; i < 50; i++ {
		input[i] = entities.MysqlCardPriceHistory{
			ID:              int64(i + 1),
			Name:            "Card",
			SetName:         "SET",
			CollectorNumber: "123",
			Foil:            i%2 == 0,
			LastPrice:       float64(i + 1),
			OldPrice:        float64(i),
			PriceChange:     1.0,
			LastUpdate:      &fixedTime,
		}
	}

	result := CardPriceHistoryToCardsDomain(input)
	assert.Equal(t, 50, len(result))
	assert.Equal(t, 50, cap(result))
}

// Helper function to create float pointers
func floatPtr(f float64) *float64 {
	return &f
}
