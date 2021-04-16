package main

import (
	"time"
	 "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"net/http"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"encoding/json"
	"bytes"
)

type Item struct {
    Year   string `json:"Year"`
    Month  string `json:"Month"`
    Cupcakes   string `json:Cupcakes"`
    ModifiedTime string `json:ModifiedTime"`
    Event string `json:"Event"`
    Id string `json:"Id"`
}


type DynamoDBEvent struct {
    Records []DynamoDBEventRecord `json:"Records"`
}

type DynamoDBEventRecord struct {
    AWSRegion      string                       `json:"awsRegion"`
    Change         DynamoDBStreamRecord         `json:"dynamodb"`
    EventID        string                       `json:"eventID"`
    EventName      string                       `json:"eventName"`
    EventSource    string                       `json:"eventSource"`
    EventVersion   string                       `json:"eventVersion"`
    EventSourceArn string                       `json:"eventSourceARN"`
    UserIdentity   *events.DynamoDBUserIdentity `json:"userIdentity,omitempty"`
}

type DynamoDBStreamRecord struct {
    ApproximateCreationDateTime events.SecondsEpochTime             `json:"ApproximateCreationDateTime,omitempty"`
    // changed to map[string]*dynamodb.AttributeValue
    Keys                        map[string]*dynamodb.AttributeValue `json:"Keys,omitempty"`
    // changed to map[string]*dynamodb.AttributeValue
    NewImage                    map[string]*dynamodb.AttributeValue `json:"NewImage,omitempty"`
    // changed to map[string]*dynamodb.AttributeValue
    OldImage                    map[string]*dynamodb.AttributeValue `json:"OldImage,omitempty"`
    SequenceNumber              string                              `json:"SequenceNumber"`
    SizeBytes                   int64                               `json:"SizeBytes"`
    StreamViewType              string                              `json:"StreamViewType"`
}



func main() {
        lambda.Start(HandleRequest)
}

    
func HandleRequest(event DynamoDBEvent) error {

  endpoint := "https://search-my-project-2-biohsm5dz4d6hd5ggjyhplgqwm.ap-south-1.es.amazonaws.com/border_free_new/"

  credentials := credentials.NewEnvCredentials()
  signer := v4.NewSigner(credentials)
  
    client := &http.Client{}
  
	for _ , record := range event.Records{
		change := record.Change
		newImage := change.NewImage
		
		var newItem Item
	  
	    err := dynamodbattribute.UnmarshalMap(newImage, &newItem)
	    if err != nil {
	        return err
	    }		
		
		eventName := record.EventName
		newItem.Event = eventName

		var jsonData []byte 
		jsonData , err= json.Marshal(newItem)
		
		if err != nil{
			return err
		}
		
		jsonDataNew := bytes.NewReader(jsonData)

		if eventName == "REMOVE"{
		
			deleteEp := endpoint+"_doc/"+newItem.Id
			req, err := http.NewRequest(http.MethodDelete, deleteEp, nil)
			if err!=nil{
				return err
			}

		  //req.Header.Add("Content-Type", "application/json")
		  region := "ap-south-1" // e.g. us-east-1
		  service := "es"
		  // Sign the request, send it, and print the response
		  signer.Sign(req, nil, service, region, time.Now())	
			_, err = client.Do(req)
			if err != nil{
				return err
			}  		
		
		}else{
		
			insertEp := endpoint+"_doc/"+newItem.Id
			req, err := http.NewRequest(http.MethodPost, insertEp, jsonDataNew)
			if err!=nil{
				return err
			}
		
		  req.Header.Add("Content-Type", "application/json")
		  region := "ap-south-1" // e.g. us-east-1
		  service := "es"
		  // Sign the request, send it, and print the response
		  signer.Sign(req, jsonDataNew, service, region, time.Now())	
			_, err = client.Do(req)
			if err != nil{
				return err
			}  		
		
		}
	}

	return nil    		
}




