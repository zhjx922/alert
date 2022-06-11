package input

import (
	"bufio"
	"errors"
	"io"
	"os"
	"time"
)

type Reader interface {
	Read() ([]byte,error)
	End()
}

type File struct {
	file *os.File
	name   string
	reader *bufio.Reader
	count  int
	done   chan struct{}
}

var ErrorDone = errors.New("DONE")

func NewFile(name string) (*File, error) {
	file, err := os.Open(name)

	if err != nil {
		return nil, err
	}

	file.Seek(0, os.SEEK_END)

	reader := bufio.NewReader(file)

	return &File{
		file: file,
		reader: reader,
		done: make(chan struct{}),
	}, nil
}

func (f *File) Read() ([]byte, error) {
	for {
		select {
		case <-f.done:
			return nil, ErrorDone
		default:
			content, err := f.reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(500 * time.Millisecond)
				}

				break
			}

			return content, nil
		}
	}

	return nil, ErrorDone
}

func (f *File) End() {
	f.file.Close()
	close(f.done)
}
