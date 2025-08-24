package ingress

type Repository struct {
	Auth       AuthServicePorts
	Handler    HandlerPorts
	Health     HealthServicePorts
	Role       RoleServicePorts
	Token      TokenServicePorts
	User       UserServicePorts
	Permission PermissionServicePorts
}
