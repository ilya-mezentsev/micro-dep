package shared

import (
	"net/http"
)

func Get(opts Opts) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, opts.Address, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range opts.Headers {
		req.Header.Set(key, val)
	}

	client := http.Client{
		Timeout: opts.Timeout,
	}

	return client.Do(req)
}
