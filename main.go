package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OldPanda/go-isbn"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type requestParams struct {
	QueryStringParameters isbnParam `json:"queryStringParameters"`
}

type isbnParam struct {
	ISBN10 string `json:"isbn10,omitempty"`
}

// HandleLambdaEvent ...
func HandleLambdaEvent(ctx context.Context, eventJSON json.RawMessage) (string, error) {
	var params requestParams
	if err := json.Unmarshal(eventJSON, &params); err != nil {
		return "", fmt.Errorf("Failed to parse url parameters: %v\nError: %v", string(eventJSON), err)
	}

	isbn10 := params.QueryStringParameters.ISBN10
	if isbn10 == "" {
		log.Warn("ISBN10 is not given")
		return "", nil
	}

	isbn13, err := isbn.ConvertToIsbn13(isbn10)
	if err != nil {
		log.Warn("Cannot convert given isbn10: %v to isbn13", isbn10)
		return "", nil
	}
	return fmt.Sprintf("%v", isbn13), nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
