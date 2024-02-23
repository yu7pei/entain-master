package service

import (
	"golang.org/x/net/context"
	"sport/db"
	"sport/proto/sports"
)

type Sports interface {
	// ListEvents will return a collection of sports events.
	ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error)

	// UpdateWinner return whether success update a winner to an event
	UpdateWinner(ctx context.Context, in *sports.UpdateWinnerRequest) (*sports.UpdateWinnerResponse, error)

	// GetSingleEventById return an event by ID
	GetSingleEventById(ctx context.Context, in *sports.GetSingleEventByIdRequest) (*sports.GetSingleEventIdResponse, error)
}

// sportsService implements the Sports interface.
type sportsService struct {
	sportsRepo db.SportsRepo
}

func (s *sportsService) GetSingleEventById(ctx context.Context, in *sports.GetSingleEventByIdRequest) (*sports.GetSingleEventIdResponse, error) {
	event, _, err := s.sportsRepo.GetByID(in.Id)
	if err != nil {
		return nil, err
	}

	return &sports.GetSingleEventIdResponse{Event: event}, nil
}

func (s *sportsService) UpdateWinner(ctx context.Context, in *sports.UpdateWinnerRequest) (*sports.UpdateWinnerResponse, error) {
	isSuccess, err := s.sportsRepo.UpdateWinner(in)
	if err != nil {
		return nil, err
	}

	return &sports.UpdateWinnerResponse{IsSuccess: isSuccess}, nil
}

// NewSportsService instantiates and returns a new sportsService.
func NewSportsService(sportsRepo db.SportsRepo) Sports {
	return &sportsService{sportsRepo}
}

func (s *sportsService) ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error) {
	events, err := s.sportsRepo.List(in.Filter, in.Order)
	if err != nil {
		return nil, err
	}

	return &sports.ListEventsResponse{Events: events}, nil
}
