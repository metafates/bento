package bit

type Bits interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Difference[B Bits](a, b B) B {
	return a & ^b
}

func Intersects[B Bits](a, b B) bool {
	return a&b != 0
}

func Contains[B Bits](a, b B) bool {
	return a&b == b
}

func Union[B Bits](a, b B) B {
	return a | b
}
