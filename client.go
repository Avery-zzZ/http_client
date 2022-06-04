package http_client

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var client = http.Client{}
var Resp *http.Response
var Cookies = []*http.Cookie{}

func Do_request(method, addr string, header map[string][]string, payload url.Values, cookies []*http.Cookie) error {

	body := ioutil.NopCloser(strings.NewReader((payload).Encode()))
	req, err := new_request_with_cookies(method, addr, body, cookies)

	if err != nil {
		return errors.New("http_client/client.go - Do_request:\n" + err.Error())
	}

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	Resp, err = client.Do(req)
	if err != nil {
		return errors.New("http_client/client.go - Do_request:\nhttp.Client.Do: " + err.Error())
	}

	return nil
}

func new_request_with_cookies(method, addr string, body io.ReadCloser, cookies []*http.Cookie) (*http.Request, error) {

	req, err := http.NewRequest(method, addr, body)
	if err != nil {
		return nil, errors.New("http_client/client.go - new_requst_with_cookies\nhttp.newRequest: " + err.Error())
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	return req, nil
}

func Save_cookies(keys []string) error {
	if Resp == nil {
		return errors.New("http_client/client.go - Save_cookies: empty response")
	}

	cookies_map := make(map[string]*http.Cookie)
	for _, resp_cookie := range Resp.Cookies() {
		cookies_map[resp_cookie.Name] = resp_cookie
	}

	for _, key := range keys {
		cookie, ok := cookies_map[key]
		if !ok {
			return errors.New("http_client/client.go - Save_cookies: no cookie named \"" + key + "\"")
		}
		Cookies = append(Cookies, cookie)
	}

	return nil
}
