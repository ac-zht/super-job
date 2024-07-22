package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	Url     string
	Timeout int64
	Client  *http.Client
	Req     *http.Request
}

func (h *HttpClient) Get() ResponseWrapper {
	var err error
	h.Req, err = http.NewRequest("GET", h.Url, nil)
	if err != nil {
		return h.requestError(err)
	}
	return h.request()
}

func (h *HttpClient) PostParams(param string) ResponseWrapper {
	var err error
	h.Req, err = http.NewRequest("POST", h.Url, bytes.NewBufferString(param))
	if err != nil {
		return h.requestError(err)
	}
	h.Req.Header.Set("Content-Type", "application/x-www-form-Urlencoded")
	return h.request()
}

func (h *HttpClient) PostJson(body string) ResponseWrapper {
	var err error
	h.Req, err = http.NewRequest("POST", h.Url, bytes.NewBufferString(body))
	if err != nil {
		return h.requestError(err)
	}
	h.Req.Header.Set("Content-Type", "application/json")
	return h.request()
}

func (h *HttpClient) request() ResponseWrapper {
	wrapper := ResponseWrapper{StatusCode: 0, Body: "", Header: make(http.Header)}
	if h.Timeout > 0 {
		h.Client.Timeout = time.Duration(h.Timeout) * time.Second
	}
	h.setRequestHeader()
	resp, err := h.Client.Do(h.Req)
	if err != nil {
		wrapper.Body = fmt.Sprintf("执行HTTP请求错误-%s", err.Error())
		return wrapper
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		wrapper.Body = fmt.Sprintf("读取HTTP请求返回值失败-%s", err.Error())
		return wrapper
	}
	wrapper.StatusCode = resp.StatusCode
	wrapper.Body = string(body)
	wrapper.Header = resp.Header
	return wrapper
}

func (h *HttpClient) setRequestHeader() {
	h.Req.Header.Set("User-Agent", "golang/super-job")
}

func (h *HttpClient) requestError(err error) ResponseWrapper {
	return ResponseWrapper{
		StatusCode: 0,
		Body:       fmt.Sprintf("创建HTTP请求错误-%s", err.Error()),
		Header:     make(http.Header),
	}
}

type ResponseWrapper struct {
	StatusCode int
	Body       string
	Header     http.Header
}
