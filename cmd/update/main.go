package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/romain-h/freedom-countries/internal/crawler"
	"github.com/romain-h/freedom-countries/internal/csv"
	"github.com/romain-h/freedom-countries/internal/email"
	"github.com/romain-h/freedom-countries/internal/score"
	"github.com/romain-h/freedom-countries/internal/storage"
)

func process() error {
	session := session.New()
	store := storage.New(session)
	mailer := email.New(session)

	// In parallel fetch last stored file and scrap new values from website
	var scores score.Countries
	var previousScores score.Countries

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		lastFile, err := store.GetFile("last.json")
		if err != nil {
			log.Fatal(err)
		}
		previousScores, err = score.ReadBuf(*lastFile)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		scores = crawler.ScrapData()
		score.Preprocess(scores)
	}()

	wg.Wait()
	// Compare values
	diff := score.GetDiff(previousScores, scores)

	// Nothing changed
	if len(diff) == 0 {
		fmt.Println("FUCP -- no diff")
		return nil
	}

	diffHTML, err := diff.RenderEmail(store, "email_template.html")
	if err != nil {
		fmt.Println(err)
		return err
	}
	updatedList := csv.GenerateBTListWithScore(store, scores, "countries.csv")

	rawEmail := email.Raw{
		Sender:      "no-reply@isorine.xyz",
		Recipient:   os.Getenv("FCUP_EMAIL"),
		Subject:     "Human Rights country watchlist",
		Message:     "Hello\r\nPlease see the attached file for a list of changes",
		MessageHTML: *diffHTML,
		Attachments: []email.Attachment{
			{
				FileName:    "freedom-countries.csv",
				FileContent: []byte(base64.StdEncoding.EncodeToString(updatedList.Bytes())),
				ContentType: "text/csv",
			},
		},
	}
	_, err = mailer.SendRaw(rawEmail)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("FUCP -- email sent")

	latestCollection, _ := json.Marshal(scores)
	store.WriteFile("last.json", latestCollection)
	return nil
}

func handler(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	err := process()
	if err != nil {
		return "ErrorProcessin", err
	}
	return "OK", nil
}

func main() {
	lambda.Start(handler)
}
