package repository

type EndpointsRepository struct {
	server *server.Server
}

func NewEndpointsRepository(s *server.Server) *EndpointsRepository {
	return &EndpointsRepository{server: s}
}
