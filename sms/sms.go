// assumes you have the following environment variables setup for AWS session creation
// AWS_SDK_LOAD_CONFIG=1
// AWS_ACCESS_KEY_ID=XXXXXXXXXX
// AWS_SECRET_ACCESS_KEY=XXXXXXXX
// AWS_DEFAULT_REGION=us-east-1

package sms

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var sess *session.Session
var svc *sns.SNS
var msdns []string

// Initialise x
func Initialise(destinations []string) {
	msdns = destinations
	if len(msdns) <= 0 {
		log.Println("SNS: No mobile subscripers set")
	}
	log.Println("SNS: creating session")
	sess = session.Must(session.NewSession())
	log.Println("SNS: session created")

	svc = sns.New(sess)
	log.Println("SNS: service created")
}

func send(message string) {
	if len(msdns) <= 0 {
		return
	}

	attributes := map[string]*sns.MessageAttributeValue{
		"AWS.SNS.SMS.SenderID": &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String("ParcelDrop"),
		},
		"AWS.SNS.SMS.SMSType": &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String("Transactional"),
		},
	}

	for _, msdn := range msdns {
		params := &sns.PublishInput{
			Message:           aws.String(message),
			PhoneNumber:       aws.String(msdn),
			MessageAttributes: attributes,
		}
		resp, err := svc.Publish(params)

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
			return
		}

		// Pretty-print the response data.
		log.Println(resp)
	}
}

// SendCorrectCode x
func SendCorrectCode(code string) {
	log.Printf("SMS: correct code: %v\n", code)
	go func() { send("Door opened with code " + code) }()
}

// SendInvalidCode x
func SendInvalidCode(code string) {
	log.Printf("SMS: invalid code: %v\n", code)
	go func() { send("Invalid code entered " + code) }()
}

// SendDoorNotClosed x
func SendDoorNotClosed() {
	log.Println("SMS: door still open")
	go func() { send("Door not closed") }()
}

// SendRescindedCode x
func SendRescindedCode(digits *string) {
	log.Println("SMS: code rescinded")
	go func() { send("Code rescinded: " + *digits) }()
}

// SendUpdatedCode x
func SendUpdatedCode(name, digits *string) {
	log.Println("SMS: code updated")
	go func() { send("Code updated: " + *name + " [" + *digits + "]") }()
}
