package control

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/teabot/parceldrop/door"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/teabot/parceldrop/codebook"
)

type InstructionType string

type Instruction struct {
	InsType    InstructionType
	Digits     *string
	AccessCode *codebook.AccessCode
}

const (
	OpenDoor    InstructionType = "open"
	RescindCode InstructionType = "rescind"
	UpdateCode  InstructionType = "update"
)

// InitialiseSqs x
func InitialiseSqs(queueURL string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	ticker := time.NewTicker(10 * time.Second)
	go poll(ticker, svc, queueURL)
}

func poll(ticker *time.Ticker, svc *sqs.SQS, queueURL string) {
	for {
		select {
		case <-ticker.C:
			receive(svc, queueURL)
		}
	}
}

func receive(svc *sqs.SQS, queueURL string) {
	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(36000), // 10 hours
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		log.Println("CONTROL: Error polling SQS", err)
		return
	}

	if len(result.Messages) == 0 {
		log.Println("CONTROL: Received no messages")
		return
	}
	for _, msg := range result.Messages {
		body := aws.StringValue(msg.Body)
		log.Printf("CONTROL: Received messages: %v\n", body)

		payload, err := decode([]byte(body))
		if err != nil {
			log.Printf("CONTROL: Error decoding message body: %v, %v\n", body, err)
		} else {
			log.Printf("CONTROL: Processing payload: %v\n", payload)
			processPayload(payload)
		}
		resultDelete, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &queueURL,
			ReceiptHandle: result.Messages[0].ReceiptHandle,
		})

		if err != nil {
			fmt.Println("Delete Error", err)
		}

		fmt.Println("Message Deleted", resultDelete)
	}
}
func processPayload(payload *Instruction) {
	switch payload.InsType {
	case OpenDoor:
		openDoorInstruction(payload)
	case RescindCode:
		rescindInstruction(payload)
	case UpdateCode:
		updateInstruction(payload)
	default:
		log.Printf("CONTROL: Unknown instruction type: %v\n", payload.InsType)
	}
}

func openDoorInstruction(payload *Instruction) {
	log.Println("CONTROL: Remote override unlock")
	door.Unlock()
}

func rescindInstruction(payload *Instruction) {
	log.Println("CONTROL: Rescinding code")
	if payload.Digits != nil {
		err := codebook.Rescind(payload.Digits)
		if err != nil {
			log.Printf("CONTROL: Error saving access code: %v, %v\n", payload.AccessCode, err)
		}
	}
}

func updateInstruction(payload *Instruction) {
	log.Println("CONTROL: Updating code")
	if payload.AccessCode != nil {
		err := codebook.Update(payload.AccessCode)
		if err != nil {
			log.Printf("CONTROL: Error saving access code: %v, %v\n", payload.AccessCode, err)
		}
	}
}

func decode(data []byte) (*Instruction, error) {
	var p *Instruction
	err := json.Unmarshal(data, &p)
	log.Printf("CONTROL: Unmarshalled %v\n", p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
