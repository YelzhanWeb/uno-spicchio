package rest

import (
	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/YelzhanWeb/uno-spicchio/internal/ports"
)

type Handlers struct {
	Service domain.Service
}

func NewHandler(service domain.Service) ports.HttpHandlers {
	return &Handlers{
		Service: service,
	}
}

func (h *Handlers) PostUserHandler() {

}
