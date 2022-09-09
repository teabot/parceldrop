package control

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/teabot/parceldrop/codebook"
)

type InstructionType string

// Annoyingly, HomeAssistant SQS integration insists on an envelope
// https://www.home-assistant.io/integrations/aws/#sqs-notify-usage
type HomeAssistantNofication struct {
	Target  *string
	Message *string
}

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

var openFn func(string)

// InitialiseSqs x
func InitialiseSqs(queueURL string, overrideFn func(string)) {
	openFn = overrideFn
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	ticker := time.NewTicker(1 * time.Second)
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
	// log.Println("CONTROL: Polling")
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
		WaitTimeSeconds:     aws.Int64(20),
	})
	// log.Println("CONTROL: Polled")

	if err != nil {
		log.Println("CONTROL: Error polling SQS", err)
		return
	}

	if len(result.Messages) == 0 {
		return
	}
	for _, msg := range result.Messages {
		body := aws.StringValue(msg.Body)
		log.Printf("CONTROL: Received messages: %v\n", body)

		payload, err := decodeInstruction([]byte(body))
		if err != nil {
			payload, err := decodeHomeAssistantNofication([]byte(body))
			if err != nil {
				log.Printf("CONTROL: Error decoding message body: %v, %v\n", body, err)
			} else {
				processPayload(payload)
			}
		} else {
			processPayload(payload)
		}
		_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &queueURL,
			ReceiptHandle: result.Messages[0].ReceiptHandle,
		})

		if err != nil {
			fmt.Println("CONTROL: Delete Error", err)
		}
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
	openFn("control")
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

func decodeInstruction(data []byte) (*Instruction, error) {
	var p *Instruction
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func decodeHomeAssistantNofication(data []byte) (*Instruction, error) {
	var m *HomeAssistantNofication
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return decodeInstruction([]byte(*m.Message))
}
