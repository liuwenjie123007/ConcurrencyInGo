package queuing

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"testing"
)

func BenchmarkUnBufferedWrite(b *testing.B) {
	performWrite(b, temFileOrFatal())
}

func BenchmarkBufferedWrite(b *testing.B) {
	bufferedFile := bufio.NewWriter(temFileOrFatal())
	performWrite(b, bufio.NewWriter(bufferedFile))
}

func temFileOrFatal() io.Writer {
	file, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Fatal(fmt.Sprintf("error: %v", err))
	}
	return file
}

func performWrite(b *testing.B, writer io.Writer) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for bt := range take(done, repeat(done, byte(0)), b.N) {
		writer.Write([]byte{bt.(byte)})
	}
}

func take(
	done <-chan interface{},
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

func repeat(
	done <-chan interface{},
	values ...interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}
