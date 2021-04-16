package main
import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)
type cupcakeFreq struct {
	Year     string 
	Month    string
	Cupcakes string
	Id string
}
func readCsvFile(filePath string) [][]string {
    f, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal("Unable to parse file as CSV for " + filePath, err)
    }

    return records
}

func insertTable(rows [][]string) {
	sess, err1 := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	})

	if err1 != nil {
		fmt.Println(err1)
	}

	svc := dynamodb.New(sess)

	tableName := "TableCupcakes"
	var a cupcakeFreq
	for i := 0; i < len(rows); i++ {
		yearMonth := strings.Split(rows[i][0], "-")
		cupcake := rows[i][1]
		year := yearMonth[0]
		month := yearMonth[1]
		id := year+month

		a = cupcakeFreq{}
		a.Year = year
		a.Month = month
		a.Cupcakes = cupcake
		a.Id = id


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


func main() {
    records := readCsvFile("./multiTimeline.csv")
    insertTable(records)
}
