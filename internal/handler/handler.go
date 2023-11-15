package handler

type Handler struct {
	ServicePath string
}

func NewHandler(servicePath string) (h *Handler) {
	h = &Handler{
		ServicePath: servicePath,
	}

	return h
}
