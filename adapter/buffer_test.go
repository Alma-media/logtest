package adapter

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Buffer_Append(t *testing.T) {
	Convey("Given a buffer with limited capacity", t, func() {
		var (
			input  = []string{"foo\n", "bar\n", "baz\n", "man\n"}
			length = len(input) - 2
			buff   = NewBuffer(length)
		)

		Convey("check if earlier entries are replaced with the latest onces", func() {
			for _, line := range input {
				buff.Append([]byte(line))
			}

			if actual := buff.Len(); actual != length {
				t.Errorf("buffer was expected to contain %d entries", length)
			}

			if actual, expected := buff.String(), strings.Join(input[2:], ""); actual != expected {
				t.Errorf("unexpected output %q was expected to be %q", actual, expected)
			}
		})
	})
}

func Test_Buffer_Flush(t *testing.T) {
	Convey("Given a write buffer and io.Writer", t, func() {
		var (
			input  = []string{"foo\n", "bar\n", "baz\n"}
			buff   = NewBuffer(len(input))
			writer bytes.Buffer
		)

		for _, line := range input {
			buff.Append([]byte(line))
		}

		if actual, expected := buff.Len(), len(input); actual != expected {
			t.Fatalf("unexpected buffer length %d, expected length is %d", actual, expected)
		}

		Convey("check if buffer is empty when successfully flushed", func() {
			output := strings.Join(input, "")

			n, err := buff.Flush(&writer)
			if err != nil {
				t.Fatalf("unexpected error while flushing the buffer: %s", err)
			}

			if expected := len(output); n != expected {
				t.Errorf("the number of written bytes %d was expected to be %d", n, expected)
			}

			if actual := buff.Len(); actual != 0 {
				// nothing should stay in the buffer since it was fully flushed
				t.Errorf("the buffer length %d was expected to be %d", actual, 0)
			}

			if actual := writer.String(); actual != output {
				t.Errorf("unexpected output %q was expected to be %q", actual, output)
			}
		})

		Convey("check leftovers when the buffer is partially flushed", func() {
			var (
				failOn  = len(input) - 1
				wrapper = NewFailingWriter(&writer, failOn)
				output  = strings.Join(input[:failOn], "")
			)

			n, err := buff.Flush(wrapper)
			if err == nil {
				t.Fatal("error was expected")
			}

			if expected := len(output); n != expected {
				t.Errorf("the number of written bytes %d was expected to be %d", n, expected)
			}

			if actual := buff.Len(); actual != 1 {
				// last element should stay in the buffer
				t.Errorf("the buffer length %d was expected to be %d", actual, 1)
			}

			if actual := writer.String(); actual != output {
				// the content must be partially written to the writer
				t.Errorf("unexpected output %q was expected to be %q", actual, output)
			}
		})
	})
}
