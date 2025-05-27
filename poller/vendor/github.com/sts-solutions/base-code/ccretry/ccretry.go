package ccretry

import (
	"encoding/json"
	"reflect"
	"time"

	"emperror.dev/errors"
)

type Retry struct {
	maxAttempts      int
	sleep            time.Duration
	fn               func() error
	errorsToRetry    map[string]struct{}
	errorsToNotRetry map[string]struct{}
	retryCondition   func(error) bool
}

type RetryResponse struct {
	maxAttempts int
	success     bool
	attempts    map[int]*string
}

func (r RetryResponse) String() string {
	aux := struct {
		MaxAttempts int             `json:"max_attempts"`
		Success     bool            `json:"success"`
		Attempts    map[int]*string `json:"attempts"`
	}{
		MaxAttempts: r.maxAttempts,
		Success:     r.success,
		Attempts:    r.attempts,
	}

	jsonBytes, _ := json.Marshal(aux)
	return string(jsonBytes)
}

func (r RetryResponse) Attempts() map[int]*string {
	return r.attempts
}

func (r RetryResponse) NumberOfAttempts() int {
	return len(r.attempts)
}

// NewRetry creates a new Retry instance with 1 attempt and 0ms sleep
func NewRetry(fn func() error) *Retry {
	return &Retry{
		fn:               fn,
		maxAttempts:      1,
		sleep:            0,
		errorsToRetry:    make(map[string]struct{}),
		errorsToNotRetry: make(map[string]struct{}),
		retryCondition:   func(err error) bool { return true },
	}
}

// WithMaxAttempts sets the max attempts
func (r *Retry) WithMaxAttempts(attempts int) *Retry {
	r.maxAttempts = attempts
	return r
}

// WithSleep sets the sleep duration between attempts
func (r *Retry) WithSleep(sleep time.Duration) *Retry {
	r.sleep = sleep
	return r
}

// WithRetryCondition sets a custom condition function to determine if an error should be retried
func (r *Retry) WithRetryCondition(fn func(error) bool) *Retry {
	r.retryCondition = fn
	return r
}

// WithNotRetryableErrorTypes sets the error types that should not be retried
// If not specified, all errors will be retried
func (r *Retry) WithNotRetryableErrorTypes(types ...any) *Retry {
	for _, t := range types {
		r.errorsToNotRetry[reflect.TypeOf(t).String()] = struct{}{}
	}
	return r
}

// WithRetryableErrorTypes sets the error types that should be retried
// If not specified, all errors will be retried
func (r *Retry) WithRetryableErrorTypes(types ...any) *Retry {
	for _, t := range types {
		r.errorsToRetry[reflect.TypeOf(t).String()] = struct{}{}
	}
	return r
}

// Run runs the function with the configured max attempts and sleep duration
func (r *Retry) Run() (resp RetryResponse, err error) {
	if r.maxAttempts <= 0 {
		return resp, errors.New("max attempts must be greater than 0")
	}

	resp = RetryResponse{
		maxAttempts: r.maxAttempts,
		attempts:    make(map[int]*string),
	}

	for attempt := 0; attempt < r.maxAttempts; attempt++ {
		currentAttempt := attempt + 1
		resp.attempts[currentAttempt] = nil
		resp.success = true

		if err = r.fn(); err == nil {
			return resp, nil
		}

		msg := err.Error()
		resp.attempts[currentAttempt] = &msg
		resp.success = false

		if !r.shouldRetry(err) {
			break
		}

		if attempt < r.maxAttempts-1 {
			time.Sleep(r.sleep)
		}
	}

	return resp, err
}

func (r *Retry) shouldRetry(err error) bool {
	if r.retryCondition != nil && !r.retryCondition(err) {
		return false
	}

	errType := reflect.TypeOf(errors.Cause(err)).String()

	// Check if error is in the do-not-retry list
	if len(r.errorsToNotRetry) > 0 {
		if _, exists := r.errorsToNotRetry[errType]; exists {
			return false
		}
	}

	// If retry list is empty, retry all errors not in do-not-retry list
	if len(r.errorsToRetry) == 0 {
		return true
	}

	// Only retry if error is in the retry list
	_, needsRetry := r.errorsToRetry[errType]
	return needsRetry
}
