package git

type Resolver interface {
	Resolve(username, repository, filename string) (string, error)
}
