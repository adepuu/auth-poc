package ports

type HttpHandler interface {
	Run(addres string)
	SetupRoutes()
	Shutdown()
}
