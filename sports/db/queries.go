package db

const (
	eventsList = "list"
)

func getEventsQueries() map[string]string {
	return map[string]string{
		eventsList: `
			SELECT 
				id, 
				name, 
				player_one, 
				player_two, 
				arena, 
				visible, 
				winner, 
				advertised_start_time 
			FROM events
		`,
	}
}
