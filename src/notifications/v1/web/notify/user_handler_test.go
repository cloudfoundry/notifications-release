package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"

	"github.com/cloudfoundry/notifications-release/src/notifications/v81/testing/mocks"
	"github.com/cloudfoundry/notifications-release/src/notifications/v81/v1/web/notify"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotifyUser", func() {
	Context("Execute", func() {
		var (
			handler     notify.UserHandler
			writer      *httptest.ResponseRecorder
			request     *http.Request
			notifyObj   *mocks.Notify
			context     stack.Context
			connection  *mocks.Connection
			strategy    *mocks.Strategy
			errorWriter *mocks.ErrorWriter
		)

		BeforeEach(func() {
			writer = httptest.NewRecorder()
			request = &http.Request{URL: &url.URL{Path: "/users/user-123"}}
			strategy = mocks.NewStrategy()
			errorWriter = mocks.NewErrorWriter()

			database := mocks.NewDatabase()
			connection = mocks.NewConnection()
			database.ConnectionCall.Returns.Connection = connection

			context = stack.NewContext()
			context.Set("database", database)
			context.Set(notify.VCAPRequestIDKey, "some-request-id")

			notifyObj = mocks.NewNotify()
			handler = notify.NewUserHandler(notifyObj, errorWriter, strategy)
		})

		Context("when notifyObj.Execute returns a successful response", func() {
			It("returns the JSON representation of the response", func() {
				notifyObj.ExecuteCall.Returns.Response = []byte("whut")

				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(writer.Body.String()).To(Equal("whut"))
			})

			It("delegates to the notifyObj object with the correct arguments", func() {
				handler.ServeHTTP(writer, request, context)

				Expect(reflect.ValueOf(notifyObj.ExecuteCall.Receives.Connection).Pointer()).To(Equal(reflect.ValueOf(connection).Pointer()))
				Expect(notifyObj.ExecuteCall.Receives.Request).To(Equal(request))
				Expect(notifyObj.ExecuteCall.Receives.Context).To(Equal(context))
				Expect(notifyObj.ExecuteCall.Receives.GUID).To(Equal("user-123"))
				Expect(notifyObj.ExecuteCall.Receives.Strategy).To(Equal(strategy))
				Expect(notifyObj.ExecuteCall.Receives.Validator).To(BeAssignableToTypeOf(notify.GUIDValidator{}))
				Expect(notifyObj.ExecuteCall.Receives.VCAPRequestID).To(Equal("some-request-id"))
			})
		})

		Context("when notifyObj.Execute returns an error", func() {
			It("propagates the error", func() {
				notifyObj.ExecuteCall.Returns.Error = errors.New("BOOM!")
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(notifyObj.ExecuteCall.Returns.Error))
			})
		})
	})
})
