package payment

import (
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type CreateSourceFunc func(amount int64, os, typePayment, currency, txnId, userId string) (*omise.Source, error)

func CreateSource(client *omise.Client) CreateSourceFunc {
	return func(amount int64, os, typePayment, currency, txnId, userId string) (*omise.Source, error) {
		result := &omise.Source{}
		err := client.Do(result, &operations.CreateSource{
			Type:         typePayment,
			Amount:       amount,
			Currency:     currency,
			PlatformType: os,
			StoreID:      txnId,
			Name:         userId,
		})

		if err != nil {
			return nil, err
		}

		return result, nil
	}
}
