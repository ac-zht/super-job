package service

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type HttpClient struct {
	url     string
	timeout int
	client  *http.Client
	req     *http.Request
}

func (h *HttpClient) Get() ResponseWrapper {
	var err error
	h.req, err = http.NewRequest("GET", h.url, nil)
	if err != nil {
		return h.requestError(err)
	}
	return h.request()
}

func (h *HttpClient) PostParams(param string) ResponseWrapper {
	var err error
	h.req, err = http.NewRequest("POST", h.url, bytes.NewBufferString(param))
	if err != nil {
		return h.requestError(err)
	}
	h.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return h.request()
}

func (h *HttpClient) PostJson(body string) ResponseWrapper {
	var err error
	h.req, err = http.NewRequest("POST", h.url, bytes.NewBufferString(body))
	if err != nil {
		return h.requestError(err)
	}
	h.req.Header.Set("Content-Type", "application/json")
	return h.request()
}

func (h *HttpClient) request() ResponseWrapper {
	wrapper := ResponseWrapper{Code: 0, Body: "", Header: make(http.Header)}
	if h.timeout > 0 {
		h.client.Timeout = time.Duration(h.timeout) * time.Second
	}
	h.setRequestHeader()
	resp, err := h.client.Do(h.req)
	if err != nil {
		wrapper.Body = fmt.Sprintf("执行HTTP请求错误-%s", err.Error())
		return wrapper
	}
}

func (h *HttpClient) setRequestHeader() {
	h.req.Header.Set("User-Agent", "golang/super-job")
}

func (h *HttpClient) requestError(err error) ResponseWrapper {
	return ResponseWrapper{
		Code:   0,
		Body:   fmt.Sprintf("创建HTTP请求错误-%s", err.Error()),
		Header: make(http.Header),
	}
}

type ResponseWrapper struct {
	Code   int
	Body   string
	Header http.Header
}
