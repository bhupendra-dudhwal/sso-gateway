package egress

type Repository struct {
	HttpClient   HttpClientPorts
	Role         RoleRepositoryPorts
	User         UserRepositoryPorts
	LoginHistory LoginHistoryPorts
}
