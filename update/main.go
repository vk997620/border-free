package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"	
)

type cupcakeFreq struct {
	Year     string 
	Month    string
	Cupcakes string
	Id string
}
func main() {
	lambda.Start(Handler)
}

func Handler() {
	rand.Seed(time.Now().UnixNano())
	minYear := 2004
	maxYear := 2020

	minMonth := 1
	maxMonth := 12
	sess, err1 := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	})

	if err1 != nil {
		fmt.Println(err1)
	}
	svc := dynamodb.New(sess)
	tableName := "TableCupcakes"
	var a cupcakeFreq
	var Month string
	for i := 0; i < 100; i++ {
		randYear := rand.Intn(maxYear-minYear+1) + minYear
		randMonth := rand.Intn(maxMonth-minMonth+1) + minMonth

		if randYear == 2020 && randMonth == 12 {
			randMonth = 11
		}
		
		Year := strconv.Itoa(randYear)
		
		if randMonth >= 1 && randMonth <= 9 {
			Month = "0" + strconv.Itoa(randMonth)
		
		}else{
			Month = strconv.Itoa(randMonth)
		}

		cupCake := strconv.Itoa(rand.Intn(101))
		Id := Year+Month

		a = cupcakeFreq{}
		a.Year = Year
		a.Month = Month
		a.Cupcakes = cupCake
		a.Id = Id
		
		av, err2 := dynamodbattribute.MarshalMap(a)
		if err2 != nil {
			fmt.Println("Got error marshalling new cupcake item:")
			fmt.Println(err2.Error())
			os.Exit(1)
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err2 = svc.PutItem(input)
		if err2 != nil {
			fmt.Println("Got error calling PutItem:")
			fmt.Println(err2.Error())
			os.Exit(1)
		}		

	}
}
