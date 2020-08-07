package c

func main() {
	go hello()
}

func hello() { // want "no recover"
	defer func() {
		println("hi")
	}()
	panic(nil)
}
