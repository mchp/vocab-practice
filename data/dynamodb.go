package data

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type dynamoDB struct {
	db *dynamodb.DynamoDB
}

// InitDynamoDB returns a usable instance of DynamoDB
func InitDynamoDB() (Database, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)
	return &dynamoDB{svc}, nil
}

// FetchNext get one of the least recently tested vocab/translation pair
func (d *dynamoDB) FetchNext() (*Word, error) {
	return nil, nil
}

// QueryWord fetches all the translations of a vocab and the last time the translations are tested
func (d *dynamoDB) QueryWord(vocab string) (*Word, error) {
	return nil, nil
}

func (d *dynamoDB) checkExist(vocab, translation string) (bool, error) {
	return false, nil
}

// Pass should be called when the user correctly identified a vocab-translation pair
func (d *dynamoDB) Pass(vocab, translation string) error {
	return nil
}

// Input submits a new vocab translation pair into the database
func (d *dynamoDB) Input(vocab, translation string) error {
	return nil
}

// List returns all the vocab and translations in storage
func (d *dynamoDB) List() ([]*Word, error) {
	return nil, nil
}
