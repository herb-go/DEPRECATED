package member

type Directive interface {
	Execute(*Service) error
}

// DirectiveFactory member Directive create factory.
type DirectiveFactory func(loader func(v interface{}) error) (Directive, error)
