package constraint

type ValidationContext interface {
	Path() string
	Root() interface{}
}
