package store

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/clients/shared"
	"github.com/ilya-mezentsev/micro-dep/shared/errs"
)

func fetch[T any](opts shared.Opts) (T, error) {
	var (
		result   T
		response *http.Response
		err      error
	)

	response, err = shared.Get(opts)
	if err != nil {
		return result, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusOK {
		var responseBody []byte
		responseBody, err = io.ReadAll(response.Body)
		if err != nil {
			return result, err
		}

		var apiResult OkResponse[T]
		err = json.Unmarshal(responseBody, &apiResult)

		result = apiResult.Data
	} else if response.StatusCode == http.StatusUnauthorized {
		err = shared.Unauthorized
	} else {
		// fixme perhaps we should try to add extra info to this case
		err = errs.Unknown
	}

	return result, err
}
