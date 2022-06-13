package rest

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

func Request(method, endpoint string, headers, queries map[string]string, body io.Reader) (code int, bodyBuffer []byte, err error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return 0, nil, err
	}

	for key, v := range headers {
		req.Header.Set(key, v)
	}

	var q url.Values
	if queries != nil {
		q = req.URL.Query()
		for key, v := range queries {
			q.Add(key, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)

	return resp.StatusCode, buf.Bytes(), nil
}
