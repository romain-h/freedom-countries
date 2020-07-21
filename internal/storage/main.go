package storage

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage interface {
	GetFile(string) (*[]byte, error)
	WriteFile(string, []byte) error
}

type store struct {
	session *session.Session
	svc     *s3.S3
}

func (s *store) WriteFile(filename string, body []byte) error {
	uploader := s3manager.NewUploaderWithClient(s.svc)
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(body),
	}

	// Perform an upload.
	if _, err := uploader.Upload(upParams); err != nil {
		return err
	}
	return nil
}

func (s *store) GetFile(filename string) (*[]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(filename),
	}

	result, err := s.svc.GetObject(input)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(result.Body)
	return &body, err
}

func New(session *session.Session) Storage {
	return &store{
		svc:     s3.New(session),
		session: session,
	}
}
