package movers

func MoveRight(pointer int) int {
	return pointer + 1
}

func MoveLeft(pointer int) int {
	if pointer == 0 {
		return pointer
	}
	return pointer - 1
}