package sendalerts

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func SendSNS() {

	fmt.Println("creating session")
	sess := session.Must(session.NewSession())
	fmt.Println("session created")

	svc := sns.New(sess)
	fmt.Println("service created")

	params := &sns.PublishInput{
		Message:     aws.String("testing 123"),
		PhoneNumber: aws.String("+11234567"),
	}
	resp, err := svc.Publish(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
