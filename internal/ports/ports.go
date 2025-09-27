package ports

type Postgres interface {
	CreateUser(userName string, password string, role string) (int, error)
}

type Service interface {
	CreateUser(userName string, password string, role string) (int, error)
}

type HttpHandlers interface {
	PostUserHandler()
}
