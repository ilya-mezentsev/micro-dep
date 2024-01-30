package shared

import "time"

type (
	Opts struct {
		Address string
		Timeout time.Duration
		Headers map[string]string
	}

	ApiErr[E any] struct {
		StatusCode int
		Error      E
	}
)
