package operators

func Add(b byte) byte {
	return b + 1
}

func Minus(b byte) byte {
	if b == 0{
		return b
	}
	
	return b - 1
}

