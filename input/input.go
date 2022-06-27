package input

import (
	"bytes"
	"github.com/zhjx922/alert/publisher"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Input struct {
	inputs *Inputs
	files map[string]*File
	done  chan struct{}
	publisher *publisher.Publisher
}

func NewInput(inputs *Inputs) *Input {
	return &Input{
		inputs: inputs,
		files: make(map[string]*File),
		done: make(chan struct{}),
	}
}

func (i *Input) read(reader Reader)  {
	for {
		b, err := reader.Read()

		if err == ErrorDone {
			return
		}

		alert := false

AlertFor:
		for _, word := range i.inputs.IncludeLines {
			if bytes.Contains(b, []byte(word)) {
				alert = true

				if len(i.inputs.ExcludeLines) > 0 {
					for _, exWord := range i.inputs.ExcludeLines {
						// 如果需要排除，不告警
						if bytes.Contains(b, []byte(exWord)) {
							alert = false
							break AlertFor
						}
					}
				}
			}
		}

		if alert {
			i.publisher.Write(b)
		}

	}
}

func (i *Input) AddFile(name string) {
	if _, ok := i.files[name]; !ok {
		if file, err := NewFile(name); err == nil {
			i.files[name] = file

			go i.read(i.files[name])

			log.Printf("Add File:%s\n", name)
		}
	}
}

func (i *Input) RemoveFile(name string) {
	if _, ok := i.files[name]; ok {
		log.Printf("Remove File:%s\n", name)
		i.files[name].End()
		delete(i.files, name)
	}
}

func (i *Input) scan()  {

	log.Println("Scanning.....")

	files := make(map[string]bool)

	for name,_ := range i.files {
		files[name] = true
	}

	for _, path := range i.inputs.Paths {
		pList, err := filepath.Glob(path)
		if err != nil {
			continue
		}

		for _, filename := range pList {

			delete(files, filename)

			fileInfo, err := os.Stat(filename)

			if err != nil {
				continue
			}

			if fileInfo.IsDir() {
				continue
			}

			i.AddFile(filename)
		}
	}

	for name,_ := range files {
		i.RemoveFile(name)
	}
}

func (i *Input) Run(publisher *publisher.Publisher) {
	if i.inputs.ScanFrequency < 1 {
		i.inputs.ScanFrequency = 10
	}
	i.publisher = publisher
	i.scan()
	for {
		select {
		case <-i.done:
			return
		case <-time.After(time.Second * time.Duration(i.inputs.ScanFrequency)):
			i.scan()
		}
	}
}
