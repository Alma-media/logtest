package adapter

import (
	"testing"
)

func Benchmark_Buffer(b *testing.B) {
	var (
		input  = []string{"foo\n", "bar\n", "baz\n", "man\n", "wow\n", "len\n", "boo\n", "lol\n", "cap\n"}
		length = int(len(input))
		buff   = NewBuffer(length)
	)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for _, entry := range input {
			buff.Append([]byte(entry))
		}

		// buff.Flush(ioutil.Discard)

		// log.Println("-----------------------------")
		// buff.Flush(os.Stdout)

	}
}
