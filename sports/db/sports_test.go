package db

import (
	"database/sql"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"sport/proto/sports"
	"strings"
	"sync"
	"testing"
	"time"
)

const _sportsTestDB = "sports.db"

func Test_sportsRepo_applyFilter(t *testing.T) {
	type fields struct {
		db   *sql.DB
		init sync.Once
	}
	type args struct {
		query  string
		filter *sports.ListEventsRequestFilter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  []interface{}
	}{
		{
			name:   "Basic query without any filter",
			fields: fields{},
			args: args{
				getEventsQueries()[eventsList],
				nil,
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events",
		},
		{
			name:   "Filter by player",
			fields: fields{},
			args: args{
				getEventsQueries()[eventsList],
				&sports.ListEventsRequestFilter{
					Player: strPtr("bob"),
				},
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events WHERE player_one = ? OR player_two = ? ",
		},
		{
			name:   "Filter by winner",
			fields: fields{},
			args: args{
				getEventsQueries()[eventsList],
				&sports.ListEventsRequestFilter{
					Winner: strPtr("bob"),
				},
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events WHERE winner = ?",
		},
		{
			name:   "Filter by visible true",
			fields: fields{},
			args: args{
				getEventsQueries()[eventsList],
				&sports.ListEventsRequestFilter{
					Visible: boolPtr(true),
				},
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events WHERE visible = true",
		},
		{
			name:   "Filter by visible false",
			fields: fields{},
			args: args{
				getEventsQueries()[eventsList],
				&sports.ListEventsRequestFilter{
					Visible: boolPtr(false),
				},
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events WHERE visible = false",
		},
		{
			name:   "Filter by arena",
			fields: fields{},
			args: args{
				getEventsQueries()[eventsList],
				&sports.ListEventsRequestFilter{
					Arena: strPtr("Melbourne"),
				},
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events WHERE arena = ?",
		},
		{
			name:   "Filter by winner and arena",
			fields: fields{},
			args: args{
				getEventsQueries()[eventsList],
				&sports.ListEventsRequestFilter{
					Arena:  strPtr("Melbourne"),
					Winner: strPtr("bob"),
				},
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events WHERE arena = ? AND winner = ?",
		},
	}

	replacer := strings.NewReplacer("\n", "", "\t", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sportsRepo{
				db:   tt.fields.db,
				init: tt.fields.init,
			}
			got, _ := s.applyFilter(tt.args.query, tt.args.filter)
			if replacer.Replace(got) != tt.want {
				t.Errorf("applyFilter() got = %v, want %v", replacer.Replace(got), tt.want)
			}
		})
	}
}

func Test_sportsRepo_applyOrderBy(t *testing.T) {
	type fields struct {
		db   *sql.DB
		init sync.Once
	}
	type args struct {
		query   string
		orderBy *sports.ListEventsRequestOrder
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "query without order_by",
			fields: fields{},
			args: args{
				getEventsQueries()[eventsList],
				nil,
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events",
		},
		{
			name:   "query by invalid field and no direction",
			fields: fields{},
			args: args{
				query: getEventsQueries()[eventsList],
				orderBy: &sports.ListEventsRequestOrder{
					Parameter: "123",
					Direction: nil,
				},
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events",
		},
		{
			name:   "query by invalid field and with direction",
			fields: fields{},
			args: args{
				query: getEventsQueries()[eventsList],
				orderBy: &sports.ListEventsRequestOrder{
					Parameter: "123",
					Direction: strPtr("DESC"),
				},
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events",
		},
		{
			name:   "query by advertised_start_time and with direction",
			fields: fields{},
			args: args{
				query: getEventsQueries()[eventsList],
				orderBy: &sports.ListEventsRequestOrder{
					Parameter: "advertised_start_time",
					Direction: strPtr("DESC"),
				},
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events ORDER BY advertised_start_time DESC",
		},
		{
			name:   "query by advertised_start_time and without direction",
			fields: fields{},
			args: args{
				query: getEventsQueries()[eventsList],
				orderBy: &sports.ListEventsRequestOrder{
					Parameter: "advertised_start_time",
					Direction: nil,
				},
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events ORDER BY advertised_start_time",
		},
		{
			name:   "query by advertised_start_time and with incorrect direction",
			fields: fields{},
			args: args{
				query: getEventsQueries()[eventsList],
				orderBy: &sports.ListEventsRequestOrder{
					Parameter: "advertised_start_time",
					Direction: strPtr("NOOO"),
				},
			},
			want: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events ORDER BY advertised_start_time",
		},
	}

	replacer := strings.NewReplacer("\n", "", "\t", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sportsRepo{
				db:   tt.fields.db,
				init: tt.fields.init,
			}
			if got := s.applyOrderBy(tt.args.query, tt.args.orderBy); replacer.Replace(got) != tt.want {
				t.Errorf("applyOrderBy() = %v, want %v", replacer.Replace(got), tt.want)
			}
		})
	}
}

func Test_sportsRepo_UpdateWinner(t *testing.T) {
	type fields struct {
		db   *sql.DB
		init sync.Once
	}
	type args struct {
		winner *sports.UpdateWinnerRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "update winner by existed ID",
			fields: fields{},
			args: args{
				winner: &sports.UpdateWinnerRequest{
					Winner: "Bria Purdy",
					Id:     29,
				},
			},
			want: true,
		},
		{
			name:   "update winner by non-existed ID",
			fields: fields{},
			args: args{
				winner: &sports.UpdateWinnerRequest{
					Winner: "Bria Purdy",
					Id:     101229,
				},
			},
			want: false,
		},
		{
			name:   "update winner by non-existed winner",
			fields: fields{},
			args: args{
				winner: &sports.UpdateWinnerRequest{
					Winner: "NON PLAYER",
					Id:     29,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sportsRepo{
				db:   setupDb(t),
				init: tt.fields.init,
			}
			got, _ := s.UpdateWinner(tt.args.winner)
			if got != tt.want {
				t.Errorf("UpdateWinner() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sportsRepo_GetByID(t *testing.T) {
	type fields struct {
		db   *sql.DB
		init sync.Once
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sports.Event
		want1   string
		wantErr bool
	}{
		{
			name:   "Get Single Event by existed ID",
			fields: fields{},
			args: args{
				29,
			},
			want1: "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events WHERE id = ?",
		},
		{
			name:   "Get Single Event by non-existed ID",
			fields: fields{},
			args: args{
				29000,
			},
			want1:   "SELECT id, name, player_one, player_two, arena, visible, winner, advertised_start_time FROM events WHERE id = ?",
			wantErr: true,
		},
	}

	replacer := strings.NewReplacer("\n", "", "\t", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &sportsRepo{
				db:   setupDb(t),
				init: tt.fields.init,
			}
			_, got1, err := r.GetByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if replacer.Replace(got1) != tt.want1 {
				t.Errorf("GetByID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_addStatus(t *testing.T) {
	var (
		futureTime = timestamppb.New(time.Now().Add(time.Hour * 24))
		pastTime   = timestamppb.New(time.Now().Add(-time.Hour * 24))
	)
	type args struct {
		events []*sports.Event
	}
	tests := []struct {
		name string
		args args
		want []*sports.Event
	}{
		{
			name: "Single event with future time",
			args: args{
				events: []*sports.Event{
					{AdvertisedStartTime: futureTime},
					{AdvertisedStartTime: pastTime},
				},
			},
			want: []*sports.Event{
				{AdvertisedStartTime: futureTime, Status: Open},
				{AdvertisedStartTime: pastTime, Status: Closed},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addStatus(tt.args.events); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

// setup DB for test
func setupDb(t *testing.T) *sql.DB {
	var (
		err error
	)

	db, err := sql.Open("sqlite3", _sportsTestDB)
	if err != nil {
		t.Fatalf("Could not open test database. %s", err)
	}

	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, name TEXT, player_one TEXT, player_two TEXT, arena TEXT, visible INTEGER, winner TEXT, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	return db
}

func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}
