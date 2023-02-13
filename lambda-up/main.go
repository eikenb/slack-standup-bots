package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	go decryptToken()
	switch strings.ToLower(os.Getenv("ENV")) {
	case "lambda", "production", "prod":
		fmt.Println("lambda.Start")
		lambda.Start(handleStart)
	default:
		http.HandleFunc("/", handleFunc)
		fmt.Println("Listening on :3000")
		http.ListenAndServe(":3000", nil)
	}
}

func handleStart(req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       commandDispatch(req.Body, realDb()),
	}
	return response, nil
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		fmt.Fprintf(w, commandDispatch(string(raw), realDb()))
	}
}
