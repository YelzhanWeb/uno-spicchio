package http

import (
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

type Handlers struct {
	Service ports.Service
}

func NewHandler(service ports.Service) ports.HttpHandlers {
	return &Handlers{
		Service: service,
	}
}

func (h *Handlers) PostUserHandler() {

}
