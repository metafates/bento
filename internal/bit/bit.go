package bit

type Bit interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Difference[B Bit](a, b B) B {
	return a & ^b
}

func Intersects[B Bit](a, b B) bool {
	return a&b != 0
}

func Contains[B Bit](a, b B) bool {
	return a&b == b
}

func Union[B Bit](a, b B) B {
	return a | b
}
