package db

import (
	"database/sql"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

const _racingTestsDB = "racing.db"

func Test_racesRepo_applyFilter(t *testing.T) {
	type fields struct {
		db   *sql.DB
		init sync.Once
	}
	type args struct {
		query  string
		filter *racing.ListRacesRequestFilter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  []interface{}
	}{
		{
			name:   "Base Case - No filters",
			fields: fields{},
			args: args{
				getRaceQueries()[racesList],
				&racing.ListRacesRequestFilter{},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races",
		},
		{
			name:   "filter single meeting ids",
			fields: fields{},
			args: args{
				getRaceQueries()[racesList],
				&racing.ListRacesRequestFilter{
					MeetingIds: []int64{5},
				},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races WHERE meeting_id IN (?)",
		},
		{
			name:   "filter multiple meeting ids",
			fields: fields{},
			args: args{
				getRaceQueries()[racesList],
				&racing.ListRacesRequestFilter{
					MeetingIds: []int64{1, 2},
				},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races WHERE meeting_id IN (?,?)",
		},
		{
			name:   "filter with visible is true",
			fields: fields{},
			args: args{
				getRaceQueries()[racesList],
				&racing.ListRacesRequestFilter{
					Visible: boolPtr(true),
				},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races WHERE visible = true",
		},
		{
			name:   "filter with visible is false",
			fields: fields{},
			args: args{
				getRaceQueries()[racesList],
				&racing.ListRacesRequestFilter{
					Visible: boolPtr(false),
				},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races WHERE visible = false",
		},
		{
			name:   "filter with visible is false and multiple meeting ids",
			fields: fields{},
			args: args{
				getRaceQueries()[racesList],
				&racing.ListRacesRequestFilter{
					MeetingIds: []int64{1, 2},
					Visible:    boolPtr(false),
				},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races WHERE meeting_id IN (?,?) AND visible = false",
		},
	}
	replacer := strings.NewReplacer("\n", "", "\t", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &racesRepo{
				db:   tt.fields.db,
				init: tt.fields.init,
			}
			got, _ := r.applyFilter(tt.args.query, tt.args.filter)
			if replacer.Replace(got) != tt.want {
				t.Errorf("applyFilter() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}

func Test_racesRepo_applyOrderBy(t *testing.T) {
	type fields struct {
		db   *sql.DB
		init sync.Once
	}
	type args struct {
		query   string
		orderBy *racing.ListRacesRequestOrderBy
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
				getRaceQueries()[racesList],
				nil,
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races",
		},
		{
			name:   "order by invalid field and no direction",
			fields: fields{},
			args: args{
				query: getRaceQueries()[racesList],
				orderBy: &racing.ListRacesRequestOrderBy{
					Parameter: "123",
					Direction: nil,
				},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races",
		},
		{
			name:   "order by invalid field and with direction",
			fields: fields{},
			args: args{
				query: getRaceQueries()[racesList],
				orderBy: &racing.ListRacesRequestOrderBy{
					Parameter: "123",
					Direction: strPtr("DESC"),
				},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races",
		},
		{
			name:   "order by advertised_start_time field and with direction",
			fields: fields{},
			args: args{
				query: getRaceQueries()[racesList],
				orderBy: &racing.ListRacesRequestOrderBy{
					Parameter: "advertised_start_time",
					Direction: strPtr("DESC"),
				},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races ORDER BY advertised_start_time DESC",
		},
		{
			name:   "order by advertised_start_time field and no direction",
			fields: fields{},
			args: args{
				query: getRaceQueries()[racesList],
				orderBy: &racing.ListRacesRequestOrderBy{
					Parameter: "advertised_start_time",
					Direction: nil,
				},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races ORDER BY advertised_start_time",
		},
		{
			name:   "order by advertised_start_time field and with 'ACS' direction",
			fields: fields{},
			args: args{
				query: getRaceQueries()[racesList],
				orderBy: &racing.ListRacesRequestOrderBy{
					Parameter: "advertised_start_time",
					Direction: strPtr("ASC"),
				},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races ORDER BY advertised_start_time ASC",
		},
		{
			name:   "order by advertised_start_time field and with incorrect direction",
			fields: fields{},
			args: args{
				query: getRaceQueries()[racesList],
				orderBy: &racing.ListRacesRequestOrderBy{
					Parameter: "advertised_start_time",
					Direction: strPtr("NOOO"),
				},
			},
			want: "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races ORDER BY advertised_start_time",
		},
	}

	replacer := strings.NewReplacer("\n", "", "\t", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &racesRepo{
				db:   tt.fields.db,
				init: tt.fields.init,
			}

			if got := r.applyOrderBy(tt.args.query, tt.args.orderBy); replacer.Replace(got) != tt.want {
				t.Errorf("applyOrderBy() = %v, want %v", got, tt.want)
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
		races []*racing.Race
	}
	tests := []struct {
		name string
		args args
		want []*racing.Race
	}{
		{
			name: "Single race with future time",
			args: args{
				races: []*racing.Race{
					{AdvertisedStartTime: futureTime},
					{AdvertisedStartTime: pastTime},
				},
			},
			want: []*racing.Race{
				{AdvertisedStartTime: futureTime, Status: Open},
				{AdvertisedStartTime: pastTime, Status: Closed},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addStatus(tt.args.races); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_racesRepo_GetByID(t *testing.T) {
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
		want    *racing.Race
		want1   string
		wantErr bool
	}{
		{
			name:    "get by non-existed id",
			fields:  fields{},
			args:    args{id: 10000},
			want1:   "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races WHERE id = ?",
			wantErr: true,
		},
		{
			name:   "get by id",
			fields: fields{},
			args:   args{id: 1},
			want1:  "SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races WHERE id = ?",
		},
	}

	replacer := strings.NewReplacer("\n", "", "\t", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &racesRepo{
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

func setupDb(t *testing.T) *sql.DB {
	var (
		err error
	)

	db, err := sql.Open("sqlite3", _racingTestsDB)
	if err != nil {
		t.Fatalf("Could not open test database. %s", err)
	}

	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS races (id INTEGER PRIMARY KEY, meeting_id INTEGER, name TEXT, number INTEGER, visible INTEGER, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	return db
}
