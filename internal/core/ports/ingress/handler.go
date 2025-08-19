package ingress

type HandlerPorts interface {
	SetHealthHandler(healthService HealthServicePorts)
}
