package publisher

import (
	"fmt"
	"github.com/zhjx922/alert/output"
	"net/http"
	"regexp"
	"strings"
)

type Publisher struct {
	http *output.Http
	message chan []byte
}

func NewPublisher(http *output.Http) *Publisher {
	return &Publisher{
		http: http,
		message: make(chan []byte),
	}
}

func (p *Publisher) Write(message []byte)  {
	p.message <- message
}

func (p *Publisher) Monitor()  {
	for  {
		select {
			case m := <-p.message:
				p.curl(string(m))
		}
	}
}

func (p *Publisher) curl(content string) {
	reg := regexp.MustCompile(`%{content}`)
	s := reg.ReplaceAllString(p.http.Body, content)

	reader := strings.NewReader(s)
	request, err := http.NewRequest(p.http.Method, p.http.Url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, header := range p.http.Headers {
		hs := strings.SplitN(strings.Trim(header, " "), " ", 2)
		request.Header.Set(hs[0], hs[1])
	}

	client := http.Client{}
	client.Do(request)
}