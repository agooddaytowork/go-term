package ansiterm

import (
	"unicode/utf8"
)

type groundState struct {
	baseState
}

func (gs groundState) Handle(b byte) (s state, e error) {
	gs.parser.context.currentChar = b
	nextState, err := gs.baseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	switch {

	case sliceContains(printables, b):
		key := [...]byte {b}
		gs.parser.context.utfChar, _ = utf8.DecodeRune(key[:])
		return gs, gs.parser.print()

	case sliceContains(executors, b):
		return gs, gs.parser.execute()
	case utf8.RuneStart(b):
		var length int
		if b&0xC0 == 0xC0 {
			length = 1
		}
		if b&0xE0 == 0xE0 {
			length = 2
		}
		if b&0xF0 == 0xF0 {
			length = 3
		}
				
		gs.parser.context.utfRuneLength = length
		gs.parser.context.utfByte = append(gs.parser.context.utfByte, b)
		return gs, nil
	}
	return gs, nil
}
