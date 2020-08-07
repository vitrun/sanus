package d

func willPanic() { // want "no defer"
	panic("I am in panic")
}

