package email

import (
	"bytes"
	"encoding/base64"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"strings"

	"github.com/pkg/errors"
)

type Raw struct {
	Sender      string
	Recipient   string
	Subject     string
	Message     string
	MessageHTML string

	Attachments []Attachment
}

type Attachment struct {
	FileName    string
	FileContent []byte // base64 format
	ContentType string // ex : image/jpeg, text/csv, application/pdf
}

func toQuotedPrintable(s string) (string, error) {
	var ac bytes.Buffer
	w := quotedprintable.NewWriter(&ac)
	_, err := w.Write([]byte(s))
	if err != nil {
		return "", errors.Wrap(err, "write")
	}
	err = w.Close()
	if err != nil {
		return "", errors.Wrap(err, "close")
	}
	return ac.String(), nil
}

func (r *Raw) BuildEmail() []byte {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	r.SetMainHeader(writer)
	r.SetBody(writer, buf)
	r.SetAttachment(writer)
	err := writer.Close()
	if err != nil {
		panic(err)
	}

	// Strip boundary line before header (doesn't work with it present)
	s := buf.String()
	if strings.Count(s, "\n") < 2 {
		panic("invalid e-mail content")
	}
	s = strings.SplitN(s, "\n", 2)[1]
	return []byte(s)
}

func (r *Raw) SetMainHeader(writer *multipart.Writer) {
	h := make(textproto.MIMEHeader)
	h.Set("Subject", r.Subject)
	h.Set("Content-Language", "en-US")
	h.Set("Content-Type", "multipart/mixed; boundary=\""+writer.Boundary()+"\"")
	h.Set("MIME-Version", "1.0")
	_, err := writer.CreatePart(h)
	if err != nil {
		panic(err)
	}
}

func (r *Raw) SetBody(writer *multipart.Writer, buf *bytes.Buffer) {
	innerWriter := multipart.NewWriter(buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("Content-Type", "multipart/alternative; boundary=\""+innerWriter.Boundary()+"\"")
	_, err := writer.CreatePart(h)
	if err != nil {
		panic(err)
	}

	r.SetBodyText(innerWriter)
	r.SetBodyHTML(innerWriter)
	err = innerWriter.Close()
	if err != nil {
		panic(err)
	}

}

func (r *Raw) SetBodyText(writer *multipart.Writer) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Transfer-Encoding", "7bit")
	h.Set("Content-Type", "text/plain; charset=utf-8")
	part, err := writer.CreatePart(h)
	if err != nil {
		panic(err)
	}
	_, err = part.Write([]byte(r.Message))
	if err != nil {
		panic(err)
	}
}
func (r *Raw) SetBodyHTML(writer *multipart.Writer) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Transfer-Encoding", "quoted-printable")
	h.Set("Content-Type", "text/html; charset=utf-8")
	part, err := writer.CreatePart(h)
	if err != nil {
		panic(err)
	}
	content, _ := toQuotedPrintable(r.MessageHTML)
	_, err = part.Write([]byte(content))
	if err != nil {
		panic(err)
	}
}

func (r *Raw) SetAttachment(writer *multipart.Writer) {
	if len(r.Attachments) == 0 {
		return
	}
	for _, attachment := range r.Attachments {
		func(obj []byte) {
			_, err := base64.StdEncoding.DecodeString(string(obj))
			if err != nil {
				panic("ATTACHMENT.FileContent not base64 format")
			}
		}(attachment.FileContent)

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", "attachment")
		h.Set("Content-Type", attachment.ContentType+"; name=\""+attachment.FileName+"\"")
		h.Set("Content-Transfer-Encoding", "base64")
		part, err := writer.CreatePart(h)
		if err != nil {
			panic(err)
		}
		_, err = part.Write(attachment.FileContent)
		if err != nil {
			panic(err)
		}
	}
}
