package payment

import (
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"time"
)

type CreateChargeFunc func(amount int64, returnURI, currency string, exp time.Time, source *omise.Source) (*omise.Charge, error)

func CreateCharge(client *omise.Client, webhookEndPoint []string) CreateChargeFunc {

	return func(amount int64, returnURI, currency string, exp time.Time, source *omise.Source) (*omise.Charge, error) {
		charge, createCharge := &omise.Charge{
			Source: source,
		}, &operations.CreateCharge{
			Amount:           amount,
			Currency:         currency,
			ReturnURI:        returnURI,
			WebhookEndpoints: webhookEndPoint,
			Source:           source.ID,
			ExpiresAt:        &exp,
		}

		if e := client.Do(charge, createCharge); e != nil {
			return nil, e
		}

		return charge, nil
	}
}
