package email

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	charset       = "UTF-8"
	region        = "AWS_REGION"
	defaultRegion = "us-east-1"
)

// Config is a representation of email configurations
type Config struct {
	Region string
}

// DefaultConfig returns the default config
func DefaultConfig() Config {
	return Config{Region: defaultRegion}
}

// SendRequest is a type for an email send request
type SendRequest struct {
	HTMLContent string
	Recipient   string
	Sender      string
	Subject     string
}

// SendEmailFunc returns functionality to send an email
type SendEmailFunc func(emailReq SendRequest) error

// SendEmail sends an email to a recipeint given HTML content
func SendEmail(c Config) SendEmailFunc {
	return func(emailReq SendRequest) error {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(c.Region)},
		)

		// Create an SES session.
		svc := ses.New(sess)

		// Assemble the email.
		input := &ses.SendEmailInput{
			Destination: &ses.Destination{
				CcAddresses: []*string{},
				ToAddresses: []*string{
					aws.String(emailReq.Recipient),
				},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Html: &ses.Content{
						Charset: aws.String(charset),
						Data:    aws.String(emailReq.HTMLContent),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(emailReq.Subject),
				},
			},
			Source: aws.String(emailReq.Sender),
		}

		// Attempt to send the email.
		result, err := svc.SendEmail(input)

		// Display error messages if they occur.
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case ses.ErrCodeMessageRejected:
					return fmt.Errorf("%v %v", ses.ErrCodeMessageRejected, aerr.Error())
				case ses.ErrCodeMailFromDomainNotVerifiedException:
					return fmt.Errorf("%v %v", ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
				case ses.ErrCodeConfigurationSetDoesNotExistException:
					return fmt.Errorf("%v %v", ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
				default:
					return aerr
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				return err
			}
		}
		fmt.Println("Email Sent to address: " + emailReq.Recipient)
		fmt.Println(result)
		return nil
	}
}
