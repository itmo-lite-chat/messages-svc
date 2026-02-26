package example_grpc_api

type API struct {
	service ExampleService
}

func NewAPI(service ExampleService) *API {
	return &API{service: service}
}
