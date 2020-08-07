package e

type S struct{}

func (s *S) foo() {
	go func() { // want "no defer"
		println("foo")
	}()
}

func (s S) bar() {
	go func() { // want "no recover"
		defer func() {

		}()
		println("bar")
	}()
}
