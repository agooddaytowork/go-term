package ansiterm

type ansiContext struct {
	currentChar byte
	paramBuffer []byte
	interBuffer []byte
	utfByte []byte
	utfChar rune
	utfRuneLength int
}
