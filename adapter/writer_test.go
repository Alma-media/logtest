package adapter

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_BufferingWriter(t *testing.T) {
	Convey("Given a buffering writer", t, func() {
		var (
			input    = []string{"foo\n", "bar\n", "baz\n"}
			buffer   = NewBuffer(len(input))
			failing  = NewFailingWriter(nil, 0)
			buffered = NewBufferingWriter(failing, buffer)
			output   = strings.Join(input, "")
		)

		Convey("test writing to the buffer", func() {
			for _, line := range input {
				if _, err := buffered.Write([]byte(line)); err != nil {
					t.Fatal("buffered writer should not throw an error")
				}
			}

			if actual := buffer.String(); actual != output {
				t.Errorf("the output %q does not match expected result %q", actual, output)
			}
		})
	})
}
