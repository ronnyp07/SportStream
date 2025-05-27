package ccnats

import (
	"errors"

	"github.com/sts-solutions/base-code/ccmsgqueue"
	"github.com/sts-solutions/base-code/ccvalidation"
)

// PublishMessageBuilder is a builder for constructing a ccmsgqueue.PublishMessage.
type PublishMessageBuilder struct {
	message *publishMessage
}

// NewPublishMessageBuilder creates a new instance of PublishMessageBuilder.
func NewPublishMessageBuilder() *PublishMessageBuilder {
	return &PublishMessageBuilder{
		message: &publishMessage{},
	}
}

// WithSubject sets the subject for the message. Subject is mandatory.
func (mb *PublishMessageBuilder) WithSubject(subject string) *PublishMessageBuilder {
	mb.message.subject = subject
	return mb
}

// WithData sets the data for the message. Data is mandatory.
func (mb *PublishMessageBuilder) WithData(data []byte) *PublishMessageBuilder {
	mb.message.data = data
	return mb
}

// WithHeaders sets the headers for the message.
func (mb *PublishMessageBuilder) WithHeaders(headers map[string][]string) *PublishMessageBuilder {
	mb.message.headers = headers
	return mb
}

// Build validates the configuration and returns the constructed Message.
func (mb *PublishMessageBuilder) Build() (ccmsgqueue.PublishMessage, error) {
	result := mb.validate()
	if result.IsFailure() {
		return nil, errors.New(result.Error())
	}
	return mb.message, nil
}

// validate checks if all required fields are set.
func (mb *PublishMessageBuilder) validate() ccvalidation.Result {
	result := ccvalidation.Result{}

	if mb.message.subject == "" {
		result.AddErrorMessage("Message subject is missing")
	}

	if mb.message.data == nil {
		result.AddErrorMessage("Message data is missing")
	}

	return result
}
