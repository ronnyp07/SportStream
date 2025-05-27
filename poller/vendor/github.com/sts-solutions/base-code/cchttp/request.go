package cchttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"emperror.dev/errors"
	"github.com/sts-solutions/base-code/cchttp/cccontenttype"
	"github.com/sts-solutions/base-code/ccvalidation"
)

var (
	errRequestHTTPClientHasNotBeenSet error = errors.New("http client has not been set")
)

type request struct {
	body                interface{}
	url                 string
	context             context.Context
	headers             map[string]string
	response            any
	responseRawBody     *string
	responseHeaders     map[string][]string
	expectedStatusCode  int
	gotStatusCode       *int
	httpClient          Client
	httpMethod          string
	bodyContentType     cccontenttype.ContentType
	responseContentType cccontenttype.ContentType

	shouldVerifyStatusCode      bool
	shouldUnmarshalResponse     bool
	shouldReturnResponseHeaders bool
	shouldReturnRawBody         bool
	isBodySet                   bool

	httpRequest *http.Request
}

func (r *request) HTTPRequest() *http.Request {
	return r.httpRequest
}

// Do executes the HTTP request and processes the response
func (r *request) Do() error {
	if err := r.validate(); err != nil {
		return errors.Wrap(err, "validating request")
	}

	resp, err := r.httpClient.Do(r.httpRequest)
	if err != nil {
		return errors.Wrap(err, "executing request")
	}

	defer resp.Body.Close()

	err = r.processResponse(resp)
	if err != nil {
		return errors.Wrap(err, "decoding response")
	}

	return nil
}

func (r *request) getHTTPRequest() (httpReq *http.Request, err error) {
	bodyReader, err := r.getBodyReader()
	if err != nil {
		return httpReq, errors.Wrap(err, "getting body reader")
	}

	httpReq, err = http.NewRequestWithContext(r.context, r.httpMethod, r.url, bodyReader)
	if err != nil {
		return httpReq, errors.Wrap(err, "creating body reader")
	}

	for headKey, headValue := range r.headers {
		httpReq.Header.Set(headKey, headValue)
	}

	return httpReq, nil
}

func (r *request) getBodyReader() (bodyReader io.Reader, err error) {

	if !r.isBodySet {
		return http.NoBody, nil
	}

	mf := r.bodyContentType.MarshalFunc()
	bodyBytes, err := mf(r.body)
	if err != nil {
		return bodyReader, errors.Wrap(err, "marshaling request boidy")
	}

	bodyReader = bytes.NewReader(bodyBytes)
	return bodyReader, nil
}

func (r *request) processResponse(resp *http.Response) (err error) {
	bodyRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "getting body reader")
	}

	if r.gotStatusCode != nil {
		*r.gotStatusCode = resp.StatusCode
	}

	if r.shouldVerifyStatusCode &&
		resp.StatusCode != r.expectedStatusCode {
		err = fmt.Errorf("expected %d, got %d "+
			"\n\tURL: %s\n\t"+
			"Request: %v \n\t"+
			"ResponseBody: %v",
			r.expectedStatusCode, resp.StatusCode, r.url,
			r.body, string(bodyRaw))

		return errors.Wrap(err, "unexpected status code")
	}

	if r.shouldReturnResponseHeaders {
		for k, v := range resp.Header {
			r.responseHeaders[k] = v
		}
	}

	if r.shouldReturnRawBody && r.responseRawBody != nil {
		*r.responseRawBody = string(bodyRaw)
	}

	if !r.shouldUnmarshalResponse {
		return nil
	}

	err = r.responseContentType.UnmarshalFunc()(bodyRaw, r.response)
	if err != nil {
		return errors.Wrap(err, "unmarshaling json response body")
	}

	return nil
}

func (r *request) validate() error {
	result := ccvalidation.Result{}

	if r.httpClient == nil {
		result.AddError(errRequestHTTPClientHasNotBeenSet)
	}
	if r.url == "" {
		return result
	}
	return nil
}
