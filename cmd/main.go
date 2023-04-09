package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	args := os.Args[1:]

	argLength := len(args)
	if argLength < 3 {
		fmt.Println("minimum of 3 arguments required, -input, S3URI and at least 1 filter")
		os.Exit(1)
	}

	if args[0] != "-input" {
		fmt.Println("no -input flag is present")
		os.Exit(1)
	}

	s3uriArr := strings.Split(args[1], "/")
	s3uriArrLen := len(s3uriArr)
	if s3uriArrLen != 4 {
		fmt.Println("invalid S3 URI")
		os.Exit(1)
	}

	// Initialize a new AWS session
	sess := session.Must(session.NewSession())

	// Create a new S3 client
	svc := s3.New(sess)

	// Specify the S3 object details
	bucket := s3uriArr[s3uriArrLen-2]
	key := s3uriArr[s3uriArrLen-1]

	// Define the S3 Select query
	query := "select * from S3Object s where"
	var where []string

	for i := 2; i < argLength; i++ {
		res := strings.Split(args[i], "=")
		switch argKey := res[0]; argKey {
		case "-with-id":
			where = append(where, fmt.Sprintf("s.id = %s", res[1]))
		case "-from-time":
			where = append(where, fmt.Sprintf("s.\"time\" >= '%s'", res[1]))
		case "-to-time":
			where = append(where, fmt.Sprintf("s.\"time\" <= '%s'", res[1]))
		case "-with-word":
			where = append(where, fmt.Sprintf("'%s' IN s.words", res[1]))
		default:
			// fmt.Printf("%s not valid.\n", argKey)
		}
	}

	if len(where) < 1 {
		fmt.Println("invalid filter flags")
		os.Exit(1)
	}

	query = fmt.Sprintf("%s %s", query, strings.Join(where, " AND "))

	// Execute the S3 Select query
	resp, err := svc.SelectObjectContent(&s3.SelectObjectContentInput{
		Bucket:         aws.String(bucket),
		Key:            aws.String(key),
		Expression:     aws.String(query),
		ExpressionType: aws.String("SQL"),
		InputSerialization: &s3.InputSerialization{
			CompressionType: aws.String("GZIP"),
			JSON: &s3.JSONInput{
				Type: aws.String("Lines"),
			},
		},
		OutputSerialization: &s3.OutputSerialization{
			JSON: &s3.JSONOutput{
				RecordDelimiter: aws.String("\n"),
			},
		},
	})
	if err != nil {
		fmt.Println("Error executing S3 Select query:", err)
		os.Exit(1)
	}
	defer resp.EventStream.Close()

	for event := range resp.EventStream.Events() {
		switch v := event.(type) {
		case *s3.RecordsEvent:
			fmt.Println(string(v.Payload))
		}
	}

}
