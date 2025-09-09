package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

type LogRecord struct {
	LogType string `json:"log_type"`
	Preview string `json:"preview"`
	Headers string `json:"headers"`
	Body    string `json:"body"`
}

type LogRoundTripper struct {
	token *oauth2.Token
	logs  []*LogRecord
}

func (t *LogRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	t.token.SetAuthHeader(req)

	l, err := t.makeRequestLogRecord(req)
	if err != nil {
		fmt.Println("req log record err ", err)
	}
	t.logs = append(t.logs, l)

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	l, err = t.makeResponseLogRecord(resp)
	if err != nil {
		fmt.Println("resp log record err ", err)
	}
	t.logs = append(t.logs, l)

	return resp, err
}

func (t *LogRoundTripper) Logs() []*LogRecord {
	return t.logs
}

func (t *LogRoundTripper) makeRequestLogRecord(req *http.Request) (*LogRecord, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close() //  must close
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	result := &LogRecord{
		LogType: "request",
		Preview: fmt.Sprintf("%s %s", req.Method, req.URL.String()),
	}
	for header, value := range req.Header {
		for _, v := range value {
			result.Headers += fmt.Sprintf("%s: %s\n", header, v)
		}
	}
	if len(bodyBytes) > 0 {
		result.Body = fmt.Sprintf("%s", string(bodyBytes))
	}
	return result, nil
}

func (t *LogRoundTripper) makeResponseLogRecord(resp *http.Response) (*LogRecord, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close() //  must close
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	result := &LogRecord{
		LogType: "response",
		Preview: fmt.Sprintf("%s", resp.Status),
	}
	for header, value := range resp.Header {
		for _, v := range value {
			result.Headers += fmt.Sprintf("%s: %s\n", header, v)
		}
	}
	if len(bodyBytes) > 0 {
		var prettyJSON bytes.Buffer
		if err = json.Indent(&prettyJSON, bodyBytes, "", "\t"); err != nil {
			return nil, err
		}
		result.Body = fmt.Sprintf("%s", string(prettyJSON.Bytes()))
	}
	return result, nil
}

func NewLogRoundTripper() *LogRoundTripper {
	return &LogRoundTripper{logs: make([]*LogRecord, 0)}
}

type HttpClient struct {
	client    *http.Client
	transport *LogRoundTripper
}

func NewHttpClient() *HttpClient {
	transport := NewLogRoundTripper()
	return &HttpClient{client: &http.Client{Transport: transport}, transport: transport}
}

func (c *HttpClient) SetToken(token *oauth2.Token) {
	c.transport.token = token
}

func (c *HttpClient) Client() *http.Client {
	return c.client
}

func (c *HttpClient) Logs() []*LogRecord {
	return c.transport.logs
}
