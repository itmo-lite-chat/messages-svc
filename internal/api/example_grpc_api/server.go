package example_grpc_api

type Server struct {
	service ExampleService
}

func NewServer(service ExampleService) *Server {
	return &Server{service: service}
}
