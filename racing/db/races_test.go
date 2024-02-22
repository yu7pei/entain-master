package db

import (
	"database/sql"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"strings"
	"sync"
	"testing"
)

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
