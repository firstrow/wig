package wig

import (
	"sync"
	"testing"

	"github.com/firstrow/wig/testutils"
	"github.com/stretchr/testify/require"
)

func TestBuildTextChangeEvent(t *testing.T) {
	source := `package wig

import "fmt"

func add(a int, b int) {
	fmt.Printf("%d", a+b)
}`

	e := NewEditor(
		testutils.Viewport,
		nil,
	)
	buf := e.BufferFindByFilePath("testfile", true)
	buf.ResetLines()
	buf.Append(source)
	require.Equal(t, source+"\n", buf.String())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		events := e.Events.Subscribe()
		wg.Done()
		msg := <-events
		msg.Wg.Done()
		event := msg.Msg.(EventTextChange)
		require.Equal(t, EventTextChange{
			Buf:   buf,
			Start: Position{Line: 4, Char: 22},
			End:   Position{Line: 4, Char: 22},
			Text:  " int",
		}, event)
	}()

	line := CursorLineByNum(buf, 4)
	TextInsert(buf, line, 22, " int")
	require.Equal(t, "func add(a int, b int) int {\n", line.Value.String())

	wg.Wait()
}

