package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zhjx922/alert/output"
	"net/http"
	"regexp"
	"strings"
)

type Publisher struct {
	http *output.Http
	regexp *regexp.Regexp
	message chan []byte
}

func NewPublisher(http *output.Http) *Publisher {
	return &Publisher{
		http: http,
		regexp: regexp.MustCompile(`%{content}`),
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
				p.curl(m)
		}
	}
}

func (p *Publisher) format(data []byte, body []byte) []byte  {
	if bytes.Index(data, []byte("[")) == 0 {
		// array
		var j []interface{}
		json.Unmarshal([]byte(p.http.Body), &j)

		p.formatArray(j, body)
		b, _ := json.Marshal(j)
		return b
	} else {
		// object
		var j map[string]interface{}
		json.Unmarshal([]byte(p.http.Body), &j)
		p.formatObject(j, body)

		b, _ := json.Marshal(j)
		return b
	}
}

func (p *Publisher) formatArray(data []interface{}, body []byte) []interface{}  {
	for k, v := range data {
		// 检查类型
		switch v.(type) {
		case string:
			// 格式转换
			data[k] = p.regexp.ReplaceAllString(v.(string), string(body))
		case interface{}:
			data[k] = p.formatObject(v, body)
		case []interface{}:
			data[k] = p.formatArray(v.([]interface{}), body)
		}
	}

	return data
}

func (p *Publisher) formatObject(object interface{}, body []byte) interface{}  {
	data := object.(map[string]interface{})
	for k, v := range data {
		// 检查类型
		switch v.(type) {
		case string:
			// 格式转换
			data[k] = p.regexp.ReplaceAllString(v.(string), string(body))
		case interface{}:
			data[k] = p.formatObject(v, body)
		case []interface{}:
			data[k] = p.formatArray(v.([]interface{}), body)
		}
	}

	return data
}

func (p *Publisher) curl(body []byte) {
	if p.http.Format == "json" {
		body = p.format([]byte(p.http.Body), body)
	} else {
		body = p.regexp.ReplaceAll([]byte(p.http.Body), body)
	}

	content := string(body)

	reader := strings.NewReader(content)
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