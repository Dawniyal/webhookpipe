package repository

type EventsRepository struct {
	server *server.Server
}

func NewEventsRepository(s *server.Server) *EventsRepository {
	return &EventsRepository{server: s}
}
