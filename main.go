package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

type response struct {
	StatusCode int     `json:"statusCode"`
	Headers    headers `json:"headers"`
	Body       string  `json:"body"`
}

type headers struct {
	ContentType string `json:"Content-Type,omitempty"`
}

// HandleLambdaEvent ...
func HandleLambdaEvent(ctx context.Context, eventJSON json.RawMessage) (response, error) {
	var params requestParams
	if err := json.Unmarshal(eventJSON, &params); err != nil {
		return response{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"error": "failed to parse url parameters: %v\nError: %v"}`, string(eventJSON), err),
			Headers: headers{
				ContentType: "application/json",
			},
		}, nil
	}

	isbn10 := params.QueryStringParameters.ISBN10
	if isbn10 == "" {
		errMsg := "isbn10 is not given"
		log.Error(errMsg)
		return response{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"error": "%s"}`, errMsg),
			Headers: headers{
				ContentType: "application/json",
			},
		}, nil
	}

	isbn13, err := isbn.ConvertToIsbn13(isbn10)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot convert given isbn10: %v to isbn13", isbn10)
		log.Warn(errMsg)
		return response{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "%s"}`, errMsg),
			Headers: headers{
				ContentType: "application/json",
			},
		}, nil
	}

	return response{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf(`{"isbn13": "%s"}`, isbn13),
		Headers: headers{
			ContentType: "application/json",
		}}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
