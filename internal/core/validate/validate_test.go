package validate

import (
	"mtg-report/internal/core/dtos"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	validator := New()
	assert.NotNil(t, validator)
}

func TestValidator_Card(t *testing.T) {
	validator := New()

	tests := []struct {
		name    string
		card    dtos.RequestInsertCard
		wantErr bool
		errMsg  string
	}{
		{
			name: "should return nil when all fields are valid",
			card: dtos.RequestInsertCard{
				Name:            "Lightning Bolt",
				SetName:         "M21",
				CollectorNumber: "123",
				Foil:            boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "should return error when name is empty",
			card: dtos.RequestInsertCard{
				Name:            "",
				SetName:         "M21",
				CollectorNumber: "123",
				Foil:            boolPtr(true),
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "should return error when collector_number is empty",
			card: dtos.RequestInsertCard{
				Name:            "Lightning Bolt",
				SetName:         "M21",
				CollectorNumber: "",
				Foil:            boolPtr(true),
			},
			wantErr: true,
			errMsg:  "collector_number is required",
		},
		{
			name: "should return error when set_name is empty",
			card: dtos.RequestInsertCard{
				Name:            "Lightning Bolt",
				SetName:         "",
				CollectorNumber: "123",
				Foil:            boolPtr(true),
			},
			wantErr: true,
			errMsg:  "set_name is required",
		},
		{
			name: "should return error when foil is nil",
			card: dtos.RequestInsertCard{
				Name:            "Lightning Bolt",
				SetName:         "M21",
				CollectorNumber: "123",
				Foil:            nil,
			},
			wantErr: true,
			errMsg:  "foil is required",
		},
		{
			name: "should return nil when foil is false",
			card: dtos.RequestInsertCard{
				Name:            "Lightning Bolt",
				SetName:         "M21",
				CollectorNumber: "123",
				Foil:            boolPtr(false),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Card(tt.card)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_CardID(t *testing.T) {
	validator := New()

	tests := []struct {
		name    string
		parts   []string
		want    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "should return id when parts are valid",
			parts:   []string{"api", "card", "123"},
			want:    "123",
			wantErr: false,
		},
		{
			name:    "should return error when parts length is not 3",
			parts:   []string{"api", "card"},
			want:    "",
			wantErr: true,
			errMsg:  "invalid url",
		},
		{
			name:    "should return error when parts length is more than 3",
			parts:   []string{"api", "card", "123", "extra"},
			want:    "",
			wantErr: true,
			errMsg:  "invalid url",
		},
		{
			name:    "should return error when id is empty",
			parts:   []string{"api", "card", ""},
			want:    "",
			wantErr: true,
			errMsg:  "id is required",
		},
		{
			name:    "should return error when id is not numeric",
			parts:   []string{"api", "card", "abc"},
			want:    "",
			wantErr: true,
			errMsg:  "invalid id",
		},
		{
			name:    "should return error when id contains special characters",
			parts:   []string{"api", "card", "12@3"},
			want:    "",
			wantErr: true,
			errMsg:  "invalid id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validator.CardID(tt.parts)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestValidator_CardName(t *testing.T) {
	validator := New()

	tests := []struct {
		name    string
		card    dtos.RequestUpdateCard
		wantErr bool
		errMsg  string
	}{
		{
			name: "should return nil when name is valid",
			card: dtos.RequestUpdateCard{
				Name: "Lightning Bolt",
			},
			wantErr: false,
		},
		{
			name: "should return error when name is empty",
			card: dtos.RequestUpdateCard{
				Name: "",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.CardName(tt.card)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_Filters(t *testing.T) {
	validator := New()

	tests := []struct {
		name            string
		setName         string
		cardName        string
		collectorNumber string
		expectedFilters map[string]string
	}{
		{
			name:            "should return empty map when all parameters are empty",
			setName:         "",
			cardName:        "",
			collectorNumber: "",
			expectedFilters: map[string]string{},
		},
		{
			name:            "should return set_name filter when only setName is provided",
			setName:         "M21",
			cardName:        "",
			collectorNumber: "",
			expectedFilters: map[string]string{"set_name": "M21"},
		},
		{
			name:            "should return name filter when only cardName is provided",
			setName:         "",
			cardName:        "Lightning Bolt",
			collectorNumber: "",
			expectedFilters: map[string]string{"name": "Lightning Bolt"},
		},
		{
			name:            "should return collector_number filter when only collectorNumber is provided",
			setName:         "",
			cardName:        "",
			collectorNumber: "123",
			expectedFilters: map[string]string{"collector_number": "123"},
		},
		{
			name:            "should return all filters when all parameters are provided",
			setName:         "M21",
			cardName:        "Lightning Bolt",
			collectorNumber: "123",
			expectedFilters: map[string]string{
				"set_name":         "M21",
				"name":             "Lightning Bolt",
				"collector_number": "123",
			},
		},
		{
			name:            "should return partial filters when some parameters are provided",
			setName:         "M21",
			cardName:        "",
			collectorNumber: "123",
			expectedFilters: map[string]string{
				"set_name":         "M21",
				"collector_number": "123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := validator.Filters(tt.setName, tt.cardName, tt.collectorNumber)
			assert.Equal(t, tt.expectedFilters, filters)
		})
	}
}

// Helper function to create bool pointers
func boolPtr(b bool) *bool {
	return &b
}
