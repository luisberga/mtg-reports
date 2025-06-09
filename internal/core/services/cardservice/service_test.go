package cardservice

import (
	"context"
	"errors"
	"mtg-report/internal/core/domain"
	"mtg-report/internal/core/dtos"
	"mtg-report/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	repoMock := mocks.NewCardsRepositoryMock()
	logMock := mocks.NewLogMock()
	commitSize := 100

	service := New(repoMock, commitSize, logMock)

	assert.NotNil(t, service)
}

func TestService_InsertCard(t *testing.T) {
	tests := []struct {
		name      string
		request   dtos.RequestInsertCard
		setupMock func(repoMock *mocks.CardsRepositoryMock)
		want      dtos.ResponseInsertCard
		wantErr   bool
	}{
		{
			name: "should insert card successfully",
			request: dtos.RequestInsertCard{
				Name:            "Lightning Bolt",
				SetName:         "M21",
				CollectorNumber: "123",
				Foil:            boolPtr(true),
			},
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				expectedCard := domain.Cards{
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
				}
				returnCard := domain.Cards{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
				}
				repoMock.On("InsertCard", mock.Anything, expectedCard).Return(returnCard, nil)
			},
			want: dtos.ResponseInsertCard{
				ID:              1,
				Name:            "Lightning Bolt",
				Set:             "M21",
				CollectorNumber: "123",
				Foil:            true,
			},
			wantErr: false,
		},
		{
			name: "should return error when repository fails",
			request: dtos.RequestInsertCard{
				Name:            "Lightning Bolt",
				SetName:         "M21",
				CollectorNumber: "123",
				Foil:            boolPtr(true),
			},
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				expectedCard := domain.Cards{
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
				}
				repoMock.On("InsertCard", mock.Anything, expectedCard).Return(domain.Cards{}, errors.New("repository error"))
			},
			want:    dtos.ResponseInsertCard{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewCardsRepositoryMock()
			logMock := mocks.NewLogMock()

			tt.setupMock(repoMock)

			service := New(repoMock, 100, logMock)
			got, err := service.InsertCard(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "service failed to insert card")
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestService_GetCardbyID(t *testing.T) {
	fixedTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		id        string
		setupMock func(repoMock *mocks.CardsRepositoryMock)
		want      dtos.ResponseCard
		wantErr   bool
	}{
		{
			name: "should get card by id successfully",
			id:   "1",
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				returnCard := domain.Cards{
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
				}
				repoMock.On("GetCardbyID", mock.Anything, "1").Return(returnCard, nil)
			},
			want: dtos.ResponseCard{
				ID:              1,
				Name:            "Lightning Bolt",
				Set:             "M21",
				CollectorNumber: "123",
				Foil:            true,
				LastPrice:       15.50,
				OldPrice:        12.00,
				PriceChange:     3.50,
				LastUpdate:      fixedTime,
			},
			wantErr: false,
		},
		{
			name: "should handle nil last update",
			id:   "1",
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				returnCard := domain.Cards{
					ID:              1,
					Name:            "Lightning Bolt",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
					CardsDetails: domain.CardsDetails{
						LastPrice:   15.50,
						OldPrice:    12.00,
						PriceChange: 3.50,
						LastUpdate:  nil,
					},
				}
				repoMock.On("GetCardbyID", mock.Anything, "1").Return(returnCard, nil)
			},
			want: dtos.ResponseCard{
				ID:              1,
				Name:            "Lightning Bolt",
				Set:             "M21",
				CollectorNumber: "123",
				Foil:            true,
				LastPrice:       15.50,
				OldPrice:        12.00,
				PriceChange:     3.50,
				LastUpdate:      time.Time{},
			},
			wantErr: false,
		},
		{
			name: "should return error when repository fails",
			id:   "1",
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				repoMock.On("GetCardbyID", mock.Anything, "1").Return(domain.Cards{}, errors.New("repository error"))
			},
			want:    dtos.ResponseCard{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewCardsRepositoryMock()
			logMock := mocks.NewLogMock()

			tt.setupMock(repoMock)

			service := New(repoMock, 100, logMock)
			got, err := service.GetCardbyID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "service failed to get card")
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestService_GetCards(t *testing.T) {
	fixedTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		filters   map[string]string
		setupMock func(repoMock *mocks.CardsRepositoryMock)
		want      []dtos.ResponseCard
		wantErr   bool
	}{
		{
			name:    "should get cards successfully",
			filters: map[string]string{"set_name": "M21"},
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				returnCards := []domain.Cards{
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
							LastUpdate:  nil,
						},
					},
				}
				repoMock.On("GetCards", mock.Anything, map[string]string{"set_name": "M21"}).Return(returnCards, nil)
			},
			want: []dtos.ResponseCard{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					Set:             "M21",
					CollectorNumber: "123",
					Foil:            true,
					LastPrice:       15.50,
					OldPrice:        12.00,
					PriceChange:     3.50,
					LastUpdate:      fixedTime,
				},
				{
					ID:              2,
					Name:            "Counterspell",
					Set:             "M21",
					CollectorNumber: "456",
					Foil:            false,
					LastPrice:       8.00,
					OldPrice:        10.00,
					PriceChange:     -2.00,
					LastUpdate:      time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name:    "should return empty slice when no cards found",
			filters: map[string]string{"set_name": "UNKNOWN"},
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				repoMock.On("GetCards", mock.Anything, map[string]string{"set_name": "UNKNOWN"}).Return([]domain.Cards{}, nil)
			},
			want:    []dtos.ResponseCard{},
			wantErr: false,
		},
		{
			name:    "should return error when repository fails",
			filters: map[string]string{"set_name": "M21"},
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				repoMock.On("GetCards", mock.Anything, map[string]string{"set_name": "M21"}).Return(nil, errors.New("repository error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewCardsRepositoryMock()
			logMock := mocks.NewLogMock()

			tt.setupMock(repoMock)

			service := New(repoMock, 100, logMock)
			got, err := service.GetCards(context.Background(), tt.filters)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "service failed to get card")
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestService_UpdateCard(t *testing.T) {
	tests := []struct {
		name      string
		request   dtos.RequestUpdateCard
		setupMock func(repoMock *mocks.CardsRepositoryMock)
		want      dtos.ResponseInsertCard
		wantErr   bool
	}{
		{
			name: "should update card successfully",
			request: dtos.RequestUpdateCard{
				ID:   "1",
				Name: "Lightning Bolt Updated",
			},
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				expectedUpdateCard := domain.UpdateCard{
					ID:   1,
					Name: "Lightning Bolt Updated",
				}
				returnCard := domain.Cards{
					ID:              1,
					Name:            "Lightning Bolt Updated",
					SetName:         "M21",
					CollectorNumber: "123",
					Foil:            true,
				}
				repoMock.On("UpdateCard", mock.Anything, expectedUpdateCard).Return(returnCard, nil)
			},
			want: dtos.ResponseInsertCard{
				ID:              1,
				Name:            "Lightning Bolt Updated",
				Set:             "M21",
				CollectorNumber: "123",
				Foil:            true,
			},
			wantErr: false,
		},
		{
			name: "should return error when id is invalid",
			request: dtos.RequestUpdateCard{
				ID:   "invalid",
				Name: "Lightning Bolt Updated",
			},
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				// No mock setup needed as parsing should fail
			},
			want:    dtos.ResponseInsertCard{},
			wantErr: true,
		},
		{
			name: "should return error when repository fails",
			request: dtos.RequestUpdateCard{
				ID:   "1",
				Name: "Lightning Bolt Updated",
			},
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				expectedUpdateCard := domain.UpdateCard{
					ID:   1,
					Name: "Lightning Bolt Updated",
				}
				repoMock.On("UpdateCard", mock.Anything, expectedUpdateCard).Return(domain.Cards{}, errors.New("repository error"))
			},
			want:    dtos.ResponseInsertCard{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewCardsRepositoryMock()
			logMock := mocks.NewLogMock()

			tt.setupMock(repoMock)

			service := New(repoMock, 100, logMock)
			got, err := service.UpdateCard(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.request.ID == "invalid" {
					assert.Contains(t, err.Error(), "service failed to parse id in update card")
				} else {
					assert.Contains(t, err.Error(), "service failed to update card")
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestService_DeleteCard(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(repoMock *mocks.CardsRepositoryMock)
		wantErr   bool
	}{
		{
			name: "should delete card successfully",
			id:   "1",
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				repoMock.On("DeleteCard", mock.Anything, "1").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should return error when repository fails",
			id:   "1",
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				repoMock.On("DeleteCard", mock.Anything, "1").Return(errors.New("repository error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewCardsRepositoryMock()
			logMock := mocks.NewLogMock()

			tt.setupMock(repoMock)

			service := New(repoMock, 100, logMock)
			err := service.DeleteCard(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "service failed to delete card")
			} else {
				assert.NoError(t, err)
			}

			repoMock.AssertExpectations(t)
		})
	}
}

func TestService_GetCardHistory(t *testing.T) {
	fixedTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		id        string
		setupMock func(repoMock *mocks.CardsRepositoryMock)
		want      []dtos.ResponseCard
		wantErr   bool
	}{
		{
			name: "should get card history successfully",
			id:   "1",
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				returnCards := []domain.Cards{
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
						ID:              1,
						Name:            "Lightning Bolt",
						SetName:         "M21",
						CollectorNumber: "123",
						Foil:            true,
						CardsDetails: domain.CardsDetails{
							LastPrice:   12.00,
							OldPrice:    10.00,
							PriceChange: 2.00,
							LastUpdate:  nil,
						},
					},
				}
				repoMock.On("GetCardHistory", mock.Anything, "1").Return(returnCards, nil)
			},
			want: []dtos.ResponseCard{
				{
					ID:              1,
					Name:            "Lightning Bolt",
					Set:             "M21",
					CollectorNumber: "123",
					Foil:            true,
					LastPrice:       15.50,
					OldPrice:        12.00,
					PriceChange:     3.50,
					LastUpdate:      fixedTime,
				},
				{
					ID:              1,
					Name:            "Lightning Bolt",
					Set:             "M21",
					CollectorNumber: "123",
					Foil:            true,
					LastPrice:       12.00,
					OldPrice:        10.00,
					PriceChange:     2.00,
					LastUpdate:      time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "should return empty slice when no history found",
			id:   "999",
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				repoMock.On("GetCardHistory", mock.Anything, "999").Return([]domain.Cards{}, nil)
			},
			want:    []dtos.ResponseCard{},
			wantErr: false,
		},
		{
			name: "should return error when repository fails",
			id:   "1",
			setupMock: func(repoMock *mocks.CardsRepositoryMock) {
				repoMock.On("GetCardHistory", mock.Anything, "1").Return(nil, errors.New("repository error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewCardsRepositoryMock()
			logMock := mocks.NewLogMock()

			tt.setupMock(repoMock)

			service := New(repoMock, 100, logMock)
			got, err := service.GetCardHistory(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "service failed to get card history")
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
			repoMock.AssertExpectations(t)
		})
	}
}

// Helper function to create bool pointers
func boolPtr(b bool) *bool {
	return &b
}
