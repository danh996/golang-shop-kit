package twilio

import (
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// MessagingClient class
type MessagingClient struct {
	client      *twilio.RestClient
	phoneNumber string
}

// NewMessagingClient constructor
func NewMessagingClient(phoneNumber, subaccountSid string) *MessagingClient {
	return &MessagingClient{
		client: twilio.NewRestClientWithParams(twilio.RestClientParams{
			AccountSid: subaccountSid,
		}),
		phoneNumber: phoneNumber,
	}
}

// Send SMS function
func (s *MessagingClient) Send(to, body string) error {
	params := &openapi.CreateMessageParams{}
	params.SetFrom(s.phoneNumber).SetTo(to).SetBody(body)
	_, err := s.client.ApiV2010.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("ApiV2010.CreateMessage: %w", err)
	}
	return nil
}
