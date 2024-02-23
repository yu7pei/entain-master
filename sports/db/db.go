package db

import (
	"time"

	"syreclabs.com/go/faker"
)

func (s *sportsRepo) seed() error {
	statement, err := s.db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, name TEXT, player_one TEXT, player_two TEXT, arena TEXT, visible INTEGER, winner TEXT, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	for i := 1; i <= 100; i++ {
		statement, err = s.db.Prepare(`INSERT OR IGNORE INTO events(id, name, player_one, player_two, arena, visible, winner, advertised_start_time) VALUES (?,?,?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Name().Name(),
				faker.Name().Name(),
				faker.Name().Name(),
				faker.Address().State(),
				faker.Number().Between(0, 1),
				faker.Name().Name(),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
			)
		}
	}

	return err
}
