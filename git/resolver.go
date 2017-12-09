package git

type Resolver interface {
	Resolve(filename string) (string, error)
}
