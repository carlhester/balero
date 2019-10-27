package sendalerts

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func SendSNS(message string, phone string) {
	fmt.Println("SendSNS start")
	sess := session.Must(session.NewSession())

	svc := sns.New(sess)

	params := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phone),
	}
	resp, err := svc.Publish(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}
