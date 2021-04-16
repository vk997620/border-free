package main

import (
	"fmt"
	"time"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	//"strings"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"io/ioutil"
)
type Response struct {
    StatusCode int               `json:"statusCode"`
    Headers    map[string]string `json:"headers"`
    Body       string            `json:"body"`
}

func main() {
	lambda.Start(Handler)
}

func Handler() (Response, error) {
  // Basic information for the Amazon Elasticsearch Service domain
  endpoint := "https://search-my-project-2-biohsm5dz4d6hd5ggjyhplgqwm.ap-south-1.es.amazonaws.com/border_free_new/_search?size=500"


  //body := strings.NewReader(json)

  // Get credentials from environment variables and create the AWS Signature Version 4 signer
  credentials := credentials.NewEnvCredentials()
  signer := v4.NewSigner(credentials)

  // An HTTP client for sending the request
  client := &http.Client{}

  // Form the HTTP request
  req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
  
  /*if(err != nil){
  	fmt.Println(err)
  	return
  }*/
  //req.Header.Add("Content-Type", "application/json")
  region := "ap-south-1" // e.g. us-east-1
  service := "es"
  // Sign the request, send it, and print the response
  signer.Sign(req, nil, service, region, time.Now())
  
  
  respe, _ := client.Do(req)
  /*if(err != nil){
  	fmt.Println(err)
  	return
  }*/
  
  ResponseBody, _ := ioutil.ReadAll(respe.Body)
  
  /*if(err != nil){
  	fmt.Println(err)
  	return
  }*/
  
  //http.ResponseWriter.Header().Set("Access-Control-Allow-Origin","*")
  
  fmt.Println("response body : " , string(ResponseBody) )
  
  return Response{
            StatusCode: 200,
            Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
            Body:       string(ResponseBody),
        },
        nil
  }







