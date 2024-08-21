package httpserver

type resp map[string]any

type HttpServer struct {
	userService     UserService
	tokenService    TokenService
	softwareService SoftwareService
	licenseService  LicenseService
}

func NewHttpServer(
	userService UserService, tokenService TokenService, softwareService SoftwareService,
	licenseService LicenseService,
) HttpServer {
	return HttpServer{
		userService:     userService,
		tokenService:    tokenService,
		softwareService: softwareService,
		licenseService:  licenseService,
	}
}
