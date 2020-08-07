package a

func a() {
	go func() { // want "no recover"
		defer func() {
		}()
		func() {
			println("world")
		}()
	}()
}
