package retrievers

func Retriever (memAlloc map[int]byte, pointer int) string {
	return string(memAlloc[pointer])
}
