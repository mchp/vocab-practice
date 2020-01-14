package data

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type dynamoDB struct {
	db *dynamodb.DynamoDB
}

const tableName = "Vocabulary"
const vocabName = "vocab"
const translationName = "translation"
const testTimeName = "lastTested"

type row struct {
	Vocab       string `json:"vocab"`
	Translation string `json:"translation"`
	LastTested  int64  `json:"lastTested"`
}

// InitDynamoDB returns a usable instance of DynamoDB
func InitDynamoDB() (Database, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-west-1"),
		Endpoint: aws.String("http://localhost:8000"),
	}))
	svc := dynamodb.New(sess)
	return &dynamoDB{svc}, nil
}

// FetchNext get one of the least recently tested vocab/translation pair
func (d *dynamoDB) FetchNext() (*Word, error) {
	filter := expression.Name(testTimeName).AttributeNotExists().Or(
		expression.Name(testTimeName).LessThan(expression.Value(time.Now().Add(-72 * time.Hour).Unix())))
	fetches := expression.NamesList(expression.Name(vocabName), expression.Name(testTimeName))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(fetches).Build()
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	result, err := d.db.Scan(params)
	if err != nil {
		return nil, err
	}

	earliest := row{
		LastTested: math.MaxInt64,
	}
	for _, i := range result.Items {
		item := row{}

		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		if earliest.LastTested > item.LastTested {
			earliest = item
		}
	}
	if earliest.Vocab == "" {
		return nil, fmt.Errorf("no eligible vocabs to fetch")
	}
	return d.QueryWord(earliest.Vocab)
}

// QueryWord fetches all the translations of a vocab and the last time the translations are tested
func (d *dynamoDB) QueryWord(vocab string) (*Word, error) {
	result, err := d.db.Query(&dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("vocab=:v"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {S: aws.String(vocab)},
		},
	})
	if err != nil {
		return nil, err
	}
	w := &Word{Vocab: vocab}
	for _, i := range result.Items {
		item := row{}

		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		w.Translations = append(w.Translations, &TranslationAndTest{item.Translation, time.Unix(item.LastTested, 0)})
	}
	return w, nil
}

func getKey(vocab, translation string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		vocabName:       {S: aws.String(vocab)},
		translationName: {S: aws.String(translation)},
	}
}

func (d *dynamoDB) checkExist(vocab, translation string) (bool, error) {
	result, err := d.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       getKey(vocab, translation),
	})
	if err != nil {
		return false, err
	}

	item := row{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return false, err
	}
	return item.Vocab == vocab && item.Translation == translation, nil
}

// Pass should be called when the user correctly identified a vocab-translation pair
func (d *dynamoDB) Pass(vocab, translation string) error {
	exist, err := d.checkExist(vocab, translation)
	if err != nil || !exist {
		return fmt.Errorf("could not find %s -> %s: %v", vocab, translation, err)
	}
	t := strconv.Itoa(int(time.Now().Unix()))
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {N: aws.String(t)},
		},
		TableName:        aws.String(tableName),
		Key:              getKey(vocab, translation),
		ReturnValues:     aws.String("NONE"),
		UpdateExpression: aws.String(fmt.Sprintf("set %s=:r", testTimeName)),
	}
	_, err = d.db.UpdateItem(input)
	return err
}

// Input submits a new vocab translation pair into the database
func (d *dynamoDB) Input(vocab, translation string) error {
	item, err := dynamodbattribute.MarshalMap(&row{
		Vocab:       vocab,
		Translation: translation,
	})
	if err != nil {
		return err
	}
	_, err = d.db.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})

	return err
}

// List returns all the vocab and translations in storage
func (d *dynamoDB) List() ([]*Word, error) {
	return nil, fmt.Errorf("unimplemented")
}
