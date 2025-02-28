package cosmosdb

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing Cosmos DB client utilities", func() {
	It("RetryOnHttpStatusOrError: Should not retry on 401 status code", func() {
		callCount := 0
		err := RetryOnHttpStatusOrError(func() error {
			callCount += 1
			if callCount != 5 {
				return &Error{StatusCode: 401, Message: "Unauthorized"}
			}
			return &Error{StatusCode: 404, Message: "Resource Not Found"}
		}, 500)
		Expect(callCount).To(Equal(1))
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("401 : Unauthorized"))
	})

	It("RetryOnHttpStatusOrError: Should retry on 500 status code", func() {
		callCount := 0
		err := RetryOnHttpStatusOrError(func() error {
			callCount += 1
			if callCount != 5 {
				return &Error{StatusCode: 500, Message: "Internal Server Error"}
			}
			return &Error{StatusCode: 503, Message: "Service Unavailable"}
		}, 500)
		Expect(callCount).To(Equal(5))
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("503 : Service Unavailable"))
	})

	It("RetryOnHttpStatusOrError: Should retry on 'http: request canceled' error message", func() {
		callCount := 0
		err := RetryOnHttpStatusOrError(func() error {
			callCount += 1
			if callCount != 5 {
				return &Error{StatusCode: 500, Message: "http: request canceled"}
			}
			return &Error{StatusCode: 503, Message: "Service Unavailable"}
		}, 504, "http: request canceled")
		Expect(callCount).To(Equal(5))
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("503 : Service Unavailable"))
	})

	It("RetryOnHttpStatusOrError: Should retry on 'http: request canceled' error message", func() {
		callCount := 0
		err := RetryOnHttpStatusOrError(func() error {
			callCount += 1
			if callCount != 5 {
				return errors.New("http: request canceled")
			}
			return errors.New("Service Unavailable")
		}, 504, "http: request canceled")
		Expect(callCount).To(Equal(5))
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("Service Unavailable"))
	})

	It("RetryOnHttpStatusOrError: Should not retry on 'Unauthorized' error message", func() {
		callCount := 0
		err := RetryOnHttpStatusOrError(func() error {
			callCount += 1
			if callCount != 5 {
				return &Error{StatusCode: 401, Message: "Unauthorized"}
			}
			return &Error{StatusCode: 500, Message: "Server Error"}
		}, 504, "http: request canceled")
		Expect(callCount).To(Equal(1))
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("401 : Unauthorized"))
	})

	It("RetryOnHttpStatusOrError: Should retry on 'http: request canceled' or 'Bad Gateway' error messages", func() {
		callCount := 0
		err := RetryOnHttpStatusOrError(func() error {
			callCount += 1
			if callCount <= 2 {
				return &Error{StatusCode: 500, Message: "http: request canceled"}
			}
			if callCount <= 4 {
				return &Error{StatusCode: 502, Message: "Bad Gateway"}
			}
			return &Error{StatusCode: 503, Message: "Service Unavailable"}
		}, 504, "http: request canceled", "Bad Gateway")
		Expect(callCount).To(Equal(5))
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("503 : Service Unavailable"))
	})

	It("RetryOnHttpStatusOrError: Should retry on 'http: request canceled' or 'Bad Gateway' error messages", func() {
		callCount := 0
		err := RetryOnHttpStatusOrError(func() error {
			callCount += 1
			if callCount <= 2 {
				return errors.New("http: request canceled")
			}
			if callCount <= 4 {
				return errors.New("Bad Gateway")
			}
			return errors.New("Service Unavailable")
		}, 504, "http: request canceled", "Bad Gateway")
		Expect(callCount).To(Equal(5))
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("Service Unavailable"))
	})
})
