package order

type ServiceGroup struct {
	OrderService
}

var ServiceGroupApp = new(ServiceGroup)
