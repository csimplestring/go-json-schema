package constraint

type Schema interface {
	Get(k string) (Keyword, bool)
	Register(Keyword) error
}

