package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type row struct {
	Vocab       string `json:"vocab"`
	Translation string `json:"translation"`
}

func main() {
	args := os.Args[1:]
	local := len(args) > 0 && args[0] == "local"
	var sess *session.Session
	if local {
		sess = session.Must(session.NewSession(&aws.Config{
			Region:   aws.String("us-west-1"),
			Endpoint: aws.String("http://localhost:8000"),
		}))
	} else {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}
	svc := dynamodb.New(sess)
	expr, err := expression.NewBuilder().WithProjection(
		expression.NamesList(expression.Name("vocab"), expression.Name("translation"))).Build()
	if err != nil {
		fmt.Println("1:", err)
		return
	}
	result, err := svc.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
		TableName:                aws.String("Vocabulary"),
	})
	if err != nil {
		fmt.Println("2:", err)
	}

	for _, item := range result.Items {
		r := row{}
		dynamodbattribute.UnmarshalMap(item, &r)
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":r": {BOOL: aws.Bool(false)},
				":s": {S: aws.String("g")},
			},
			TableName: aws.String("Vocabulary"),
			Key: map[string]*dynamodb.AttributeValue{
				"vocab":       {S: aws.String(r.Vocab)},
				"translation": {S: aws.String(r.Translation)},
			},
			ReturnValues:     aws.String("NONE"),
			UpdateExpression: aws.String("set archived=:r, globalKey=:s"),
		}
		_, err := svc.UpdateItem(input)
		if err != nil {
			fmt.Println("3:", err)
		}
	}
}
