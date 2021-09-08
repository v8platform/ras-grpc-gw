package service

type Services struct {
	Tokens  TokensService
	Clients ClientsService
	Users   UsersService
}

func NewServices(tokens tokensService, clients ClientsService, users UsersService) *Services {
	return &Services{
		Tokens:  tokens,
		Clients: clients,
		Users:   users,
	}
}
