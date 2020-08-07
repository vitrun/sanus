package b

func main() {
	go func() {  // want "no defer"
		panic("I'm paniced")
	}()
}
