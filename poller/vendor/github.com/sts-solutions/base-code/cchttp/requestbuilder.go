package cchttp

import (
	"context"
	"net/http"

	"emperror.dev/errors"
	"github.com/sts-solutions/base-code/cchttp/cccontenttype"
	"github.com/sts-solutions/base-code/ccmiddlewares/cccorrelation"
	"github.com/sts-solutions/base-code/ccvalidation"
)

var (
	errRequestBuilderURLNotBeenSet               error = errors.New("url has not been set")
	errRequestBuilderHTTPMethodNotBeenSet        error = errors.New("http method has not been set")
	errRequestBuilderMarshalFunctionNotBeenSet   error = errors.New("marshal function has not been set")
	errRequestBuilderUnmarshalFunctionNotBeenSet error = errors.New("unmarshal function has not been set")
	errRequestBuilderContextNotBeenSet           error = errors.New("context has not been set")
)

type requestBuilder struct {
	request *request
}

// NewRequestBuilder creates a new request builder
func NewRequestBuilder() *requestBuilder {
	return &requestBuilder{
		request: &request{
			context:            context.Background(),
			expectedStatusCode: -1,
			headers:            make(map[string]string),
			responseHeaders:    make(map[string][]string),
			isBodySet:          false,
			gotStatusCode:      new(int),
		},
	}
}

// Build builds the request
func (rb *requestBuilder) Build() (r *request, err error) {

	if v := rb.validateRequest(); v != nil {
		return nil, errors.Wrap(v, "validating request")
	}

	rb.request.httpRequest, err = rb.request.getHTTPRequest()
	if err != nil {
		return nil, errors.New("preparing http request: " + err.Error())
	}

	return rb.request, nil
}

// WithContext sets the context of the request
func (rb *requestBuilder) WithContext(ctx context.Context) *requestBuilder {
	rb.request.context = ctx
	return rb
}

// WithURL sets the url to be called to perform the request
func (rb *requestBuilder) WithURL(url string) *requestBuilder {
	rb.request.url = url
	return rb
}

// WithHTTPMethod sets the http method to be used in the request
func (rb *requestBuilder) WithHTTPMethod(method string) *requestBuilder {
	rb.request.httpMethod = method
	return rb
}

// WithHTTPClient sets the HTTP client to be used for the request
func (rb *requestBuilder) WithHTTPClient(client Client) *requestBuilder {
	rb.request.httpClient = client
	return rb
}

// WithHTTPMethod sets the http method to be used in the request
func (rb *requestBuilder) WithDefaultHTTPClient() *requestBuilder {
	rb.request.httpClient = http.DefaultClient
	return rb
}

// WithExpectedStatusCode sets the expected status code for the response
func (rb *requestBuilder) WithExpectedStatusCode(code int) *requestBuilder {
	rb.request.expectedStatusCode = code
	rb.request.shouldVerifyStatusCode = true
	return rb
}

// WithStatusCode gets the status code for the response
func (rb *requestBuilder) WithStatusCode(code *int) *requestBuilder {
	rb.request.gotStatusCode = code
	return rb
}

// WithHeader sets a header for the request
func (rb *requestBuilder) WithHeader(key, value string) *requestBuilder {
	if rb.request.headers == nil {
		rb.request.headers = make(map[string]string)
	}
	rb.request.headers[key] = value
	return rb
}

// WithHeaders sets multiple headers for the request
func (rb *requestBuilder) WithHeaders(headers map[string]string) *requestBuilder {
	if rb.request.headers == nil {
		rb.request.headers = make(map[string]string)
	}
	for k, v := range headers {
		rb.request.headers[k] = v
	}
	return rb
}

// WithRequestTrackerHeader adds the request tracker header to the http request
func (rb *requestBuilder) WithRequestTrackerHeader(value string) *requestBuilder {
	return rb.WithHeader(cccorrelation.Key.String(), value)
}

// WithCorrelationIDHeader adds the correlationID header to the http request
func (rb *requestBuilder) WithCorrelationIDHeader(value string) *requestBuilder {
	return rb.WithHeader(cccorrelation.Key.String(), value)
}

// WithCorrelationIDHeaderFromContext adds the correlationID header with the
// correlationID value from the context to the http request
func (rb *requestBuilder) WithCorrelationIDHeaderFromContext(
	ctx context.Context) *requestBuilder {

	return rb.WithHeader(cccorrelation.Key.String(), cccorrelation.GetCorrelationID(ctx))
}

// Helper function to get correlation ID from context
func GetCorrelationIDFromContext(ctx context.Context) string {
	if val, ok := ctx.Value("correlationID").(string); ok {
		return val
	}
	return ""
}

// WithBody sets the request body and content type
// The content type is used to determine the marshal function to use and sets
// the content type header
func (rb *requestBuilder) WithBody(body any, contentType cccontenttype.ContentType) *requestBuilder {
	rb.request.isBodySet = true
	rb.request.body = body
	rb.request.bodyContentType = contentType

	return rb.WithHeader(cccontenttype.Key.String(), contentType.Name())
}

// WithResponse sets the response type and content type
// The content type is used to determine the unmarshal function to use
func (rb *requestBuilder) WithResponse(resp any, contentType cccontenttype.ContentType) *requestBuilder {
	rb.request.response = resp
	rb.request.shouldUnmarshalResponse = true
	rb.request.responseContentType = contentType

	return rb
}

// WithResponseRawBody sets the string to store the response raw body
func (rb *requestBuilder) WithResponseRawBody(respRawBody *string) *requestBuilder {
	rb.request.responseRawBody = respRawBody
	rb.request.shouldReturnRawBody = true

	return rb
}

// WithResponseHeaders sets the map to store the response headers
func (rb *requestBuilder) WithResponseHeaders(values map[string][]string) *requestBuilder {
	rb.request.responseHeaders = values
	rb.request.shouldReturnResponseHeaders = true
	return rb
}

// validateRequest checks if all required fields are properly set
func (rb *requestBuilder) validateRequest() error {
	result := ccvalidation.Result{}

	if rb.request.context == nil {
		result.AddError(errRequestBuilderContextNotBeenSet)
	}
	if rb.request.url == "" {
		result.AddError(errRequestBuilderURLNotBeenSet)
	}
	if rb.request.httpMethod == "" {
		result.AddError(errRequestBuilderHTTPMethodNotBeenSet)
	}

	if rb.request.isBodySet &&
		rb.request.bodyContentType.MarshalFunc() == nil {
		result.AddError(errRequestBuilderMarshalFunctionNotBeenSet)
	}

	if rb.request.shouldUnmarshalResponse &&
		rb.request.responseContentType.UnmarshalFunc() == nil {
		result.AddError(errRequestBuilderUnmarshalFunctionNotBeenSet)
	}

	if result.IsFailure() {
		return result
	}
	return nil
}
