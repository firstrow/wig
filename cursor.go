package wig

import "unicode"

type Cursor struct {
	Line                 int
	Char                 int
	PreserveCharPosition int
}

func restoreCharPosition(buf *Buffer) {
	line := CursorLine(buf)
	if line == nil {
		buf.Cursor.Char = 0
		return
	}

	if len(line.Value) == 0 {
		buf.Cursor.Char = 0
		return
	}

	if buf.Cursor.PreserveCharPosition >= len(line.Value) {
		buf.Cursor.Char = len(line.Value) - 1
	} else {
		buf.Cursor.Char = buf.Cursor.PreserveCharPosition
	}
}

func CursorInc(buf *Buffer) (moved bool) {
	line := CursorLine(buf)
	if buf.Cursor.Char < len(line.Value)-1 {
		buf.Cursor.Char++
		buf.Cursor.PreserveCharPosition = buf.Cursor.Char
		return true
	}

	if line.Next() != nil {
		buf.Cursor.Char = 0
		buf.Cursor.Line++
		buf.Cursor.PreserveCharPosition = buf.Cursor.Char
		return true
	}

	return false
}

func CursorDec(buf *Buffer) (moved bool) {
	if buf.Cursor.Char > 0 {
		buf.Cursor.Char--
		buf.Cursor.PreserveCharPosition = buf.Cursor.Char
		return true
	}

	line := CursorLine(buf)
	if line.Prev() != nil {
		chLen := max(len(line.Prev().Value)-1, 0)
		buf.Cursor.Char = chLen
		buf.Cursor.PreserveCharPosition = buf.Cursor.Char
		buf.Cursor.Line--
		return true
	}

	return false
}

func CursorLine(buf *Buffer) *Element[Line] {
	num := 0
	currentLine := buf.Lines.First()
	for currentLine != nil {
		if buf.Cursor.Line == num {
			return currentLine
		}
		currentLine = currentLine.Next()
		num++
	}
	return currentLine
}

func CursorLineByNum(buf *Buffer, num int) *Element[Line] {
	i := 0
	currentLine := buf.Lines.First()
	for currentLine != nil {
		if i == num {
			return currentLine
		}
		currentLine = currentLine.Next()
		i++
	}

	return currentLine
}

func CursorNumByLine(buf *Buffer, lookie *Element[Line]) int {
	i := 0
	currentLine := buf.Lines.First()
	for currentLine != nil {
		if currentLine == lookie {
			return i
		}
		currentLine = currentLine.Next()
		i++
	}

	return 0
}

// class of char under cursor
type chClass int

const (
	chWhitespace chClass = iota
	chPunct
	chWord
)

func CursorChClass(buf *Buffer) chClass {
	line := CursorLine(buf)

	if len(line.Value) == 0 {
		return chWhitespace
	}

	chLen := buf.Cursor.Char
	if chLen > len(line.Value)-1 {
		chLen = len(line.Value) - 1
	}

	return getChClass(line.Value[chLen])
}

// Return char under the cursor.
func CursorChar(buf *Buffer) rune {
	line := CursorLine(buf)

	if line.Value.IsEmpty() {
		return -1
	}

	return line.Value[buf.Cursor.Char]
}

func getChClass(r rune) chClass {
	if unicode.IsSpace(r) {
		return chWhitespace
	}

	if r == '_' {
		return chWord
	}

	if unicode.IsPunct(r) || unicode.IsSymbol(r) {
		return chPunct
	}

	return chWord
}

