package email

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ses "github.com/aws/aws-sdk-go/service/sesv2"
)

const CharSet = "UTF-8"

type Mailer interface {
	SendRaw(Raw) (*ses.SendEmailOutput, error)
}

type email struct {
	session *session.Session
	svc     *ses.SESV2
}

func (e *email) SendRaw(raw Raw) (*ses.SendEmailOutput, error) {
	data := raw.BuildEmail()

	input := &ses.SendEmailInput{
		Content: &ses.EmailContent{
			Raw: &ses.RawMessage{
				Data: data,
			},
		},
		Destination: &ses.Destination{
			ToAddresses: aws.StringSlice([]string{raw.Recipient}),
		},
		FromEmailAddress: aws.String(raw.Sender),
	}

	result, err := e.svc.SendEmail(input)

	return result, err
}

func New(session *session.Session) Mailer {
	return &email{
		svc:     ses.New(session),
		session: session,
	}
}
