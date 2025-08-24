package ingress

type HandlerPorts interface {
	SetHealthHandler(healthService HealthServicePorts)
	SetRoleHandler(healthService RoleServicePorts)
	SetPermissionHandler(healthService PermissionServicePorts)
	SetUserHandler(userService UserServicePorts)
	SetAuthHandler(authService AuthServicePorts)
}
