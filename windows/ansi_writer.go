// +build windows

package windowsconsole

import (
	"io"
	"os"
	ansiterm "github.com/mxseba/go-term/windows/ansiterm"
	"github.com/mxseba/go-term/windows/ansiterm/winterm"
)

// ansiWriter wraps a standard output file (e.g., os.Stdout) providing ANSI sequence translation.
type ansiWriter struct {
	file           *os.File
	fd             uintptr
	infoReset      *winterm.CONSOLE_SCREEN_BUFFER_INFO
	command        []byte
	escapeSequence []byte
	inAnsiSequence bool
	parser		   *ansiterm.AnsiParser
}

// NewAnsiWriter returns an io.Writer that provides VT100 terminal emulation on top of a
// Windows console output handle.
func NewAnsiWriter(nFile int) (io.Writer, ansiterm.AnsiEventHandler) {
//	initLogger()
	file, fd := winterm.GetStdFile(nFile)
	info, err := winterm.GetConsoleScreenBufferInfo(fd)
	if err != nil {
		return nil, nil
	}
	handler := winterm.CreateWinEventHandler(fd, file)
	parser := ansiterm.CreateParser("Ground", handler)
	
	aw := &ansiWriter{
		file:           file,
		fd:             fd,
		infoReset:      info,
		command:        make([]byte, 0, ansiterm.ANSI_MAX_CMD_LENGTH),
		escapeSequence: []byte(ansiterm.KEY_ESC_CSI),
		parser:         parser,
	}
	return aw, handler
}

func (aw *ansiWriter) Fd() uintptr {
	return aw.fd
}

// Write writes len(p) bytes from p to the underlying data stream.
func (aw *ansiWriter) Write(p []byte) (total int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

//	logger.Infof("Write: % x", p)
//	logger.Infof("Write: %s", string(p))
	return aw.parser.Parse(p)
}
