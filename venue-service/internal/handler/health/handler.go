package health

type Handler struct {
	env string
}

func New(env string) *Handler {
	return &Handler{
		env: env,
	}
}
