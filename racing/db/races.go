package db

import (
	"database/sql"
	"log"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

const Open = "OPEN"
const Closed = "CLOSED"

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter, orderBy *racing.ListRacesRequestOrderBy) ([]*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter, orderBy *racing.ListRacesRequestOrderBy) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	query = r.applyOrderBy(query, orderBy)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}
	if filter.Visible != nil {
		clauses = append(clauses, "visible = "+strconv.FormatBool(filter.GetVisible()))
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		races = append(races, &race)
	}
	races = addStatus(races)

	return races, nil
}

// this function is for users order result by column
func (r *racesRepo) applyOrderBy(query string, orderBy *racing.ListRacesRequestOrderBy) string {
	// valid columns name
	columns := []string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}

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

func addStatus(races []*racing.Race) []*racing.Race {
	for _, race := range races {

		if race.AdvertisedStartTime == nil {
			race.Status = Closed
		}

		if race.AdvertisedStartTime.AsTime().After(time.Now()) {
			race.Status = Open
		} else {
			race.Status = Closed
		}
	}

	return races
}
