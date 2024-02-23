package db

import (
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"slices"
	"sport/proto/sports"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const Open = "OPEN"
const Closed = "CLOSED"

// SportsRepo provides repository access to sports.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// List will return a list of events.
	List(filter *sports.ListEventsRequestFilter, order *sports.ListEventsRequestOrder) ([]*sports.Event, error)

	// UpdateWinner will return whether update winner to an event successful
	UpdateWinner(winner *sports.UpdateWinnerRequest) (bool, error)

	// GetByID will get a single event by id
	GetByID(id int64) (*sports.Event, string, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewSportsRepo creates a new event repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sports repository dummy data.
func (s *sportsRepo) Init() error {
	var err error

	s.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy events.
		err = s.seed()
	})

	return err
}

// List will return a list of events.
func (r *sportsRepo) List(filter *sports.ListEventsRequestFilter, orderBy *sports.ListEventsRequestOrder) ([]*sports.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventsQueries()[eventsList]

	// apply filter to our query
	query, args = r.applyFilter(query, filter)

	// apply order by to our query
	query = r.applyOrderBy(query, orderBy)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanEvents(rows)
}

func (m *sportsRepo) scanEvents(rows *sql.Rows) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.Name, &event.PlayerOne, &event.PlayerTwo, &event.Arena, &event.Visible, &event.Winner, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts := timestamppb.New(advertisedStart)

		event.AdvertisedStartTime = ts

		events = append(events, &event)
	}

	events = addStatus(events)
	return events, nil
}

// UpdateWinner will return whether update winner to an event successful
func (s *sportsRepo) UpdateWinner(winner *sports.UpdateWinnerRequest) (bool, error) {
	if winner == nil {
		return false, nil
	}
	// check whether event existed or not
	event, _, er := s.GetByID(winner.GetId())
	// if not existed return false
	if event == nil || er != nil {
		return false, er
	}

	// check whether winner is this event player
	if event.PlayerOne != winner.Winner && event.PlayerTwo != winner.Winner {
		return false, status.Error(codes.NotFound, "wrong event id, player not in it")
	}

	// query for update winner
	query := "UPDATE events SET winner = " + "'" + winner.Winner + "'" + " WHERE id = " + strconv.FormatInt(winner.Id, 10)

	_, err := s.db.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *sportsRepo) applyOrderBy(query string, orderBy *sports.ListEventsRequestOrder) string {
	// valid columns name
	columns := []string{"id", "name", "player_one", "player_one", "arena", "visible", "winner", "advertised_start_time"}

	if orderBy == nil {
		return query
	}
	if ok := slices.Contains(columns, orderBy.Parameter); !ok {
		log.Print("please input correct parameter for order by. such as 'advertised_start_time'!")
		return query
	}
	query += " ORDER BY " + orderBy.GetParameter()

	// append direction after query
	if orderBy.Direction != nil {
		direction := strings.ToUpper(orderBy.GetDirection())
		if direction == "ASC" {
			query += " ASC"
		} else if direction == "DESC" {
			query += " DESC"
		}
	}
	return query
}

func (s *sportsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if filter.Player != nil {
		clauses = append(clauses, "player_one = ? OR player_two = ? ")
		args = append(args, filter.GetPlayer(), filter.GetPlayer())

	}

	if filter.Arena != nil {
		clauses = append(clauses, "arena = ?")
		args = append(args, filter.GetArena())
	}

	if filter.Winner != nil {
		clauses = append(clauses, "winner = ?")
		args = append(args, filter.GetWinner())
	}

	if filter.Visible != nil {
		clauses = append(clauses, "visible = "+strconv.FormatBool(filter.GetVisible()))
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func addStatus(events []*sports.Event) []*sports.Event {
	for _, event := range events {
		if event.AdvertisedStartTime == nil {
			event.Status = Closed
		}

		if event.AdvertisedStartTime.AsTime().After(time.Now()) {
			event.Status = Open
		} else {
			event.Status = Closed
		}
	}
	return events
}

// GetByID will get a single event by id
func (r *sportsRepo) GetByID(id int64) (*sports.Event, string, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventsQueries()[eventsList]

	query += " WHERE id = ?"
	args = append(args, id)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, query, err
	}

	res, err := r.scanEvents(rows)
	if err != nil {
		return nil, query, err
	}

	if len(res) == 0 {
		// Event was not found
		err = status.Error(codes.NotFound, "Event was not found")
		return nil, query, err
	}

	return res[0], query, nil
}
