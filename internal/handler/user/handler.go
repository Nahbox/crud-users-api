package user

type Handler struct {
	db DB
}

func NewHandler(db DB) *Handler {
	handler := &Handler{
		db: db,
	}

	return handler
}
