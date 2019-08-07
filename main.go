package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
	"gopkg.in/robfig/cron.v2"
)

var (
	currentMenu   *menu
	slackURL      = ""
	slackChannel  = "#schachcafe"
	slackUsername = ":hamburger:"
)

func main() {
	flag.StringVar(&slackURL, "slack-url", slackURL, "URL of the incoming webhook to use (required)")
	flag.StringVar(&slackChannel, "slack-channel", slackChannel, "Slack channel to send the message to")
	flag.StringVar(&slackUsername, "slack-username", slackUsername, "Username to use when sending messages")
	flag.Parse()

	if slackURL == "" {
		log.Fatal("-slack-url cannot be left empty")
	}

	c := cron.New()
	wg := sync.WaitGroup{}
	wg.Add(1)

	// at 1pm reset everything
	_, err := c.AddFunc("0 13 * * MON-FRI", func() {
		log.Println("Resetting state.")
		currentMenu = nil
	})
	if err != nil {
		log.Fatalf("Failed to setup cron handler: %v", err)
	}

	// at 11.30 start attempt to report this day's menu for up to 45 minutes
	_, err = c.AddFunc("30 11 * * MON-FRI", func() {
		now := time.Now()
		delay := 1 * time.Minute
		timeout := now.Add(45 * time.Minute)

		for now.Before(timeout) {
			if notifier(now) {
				break
			}

			time.Sleep(delay)
			now = time.Now()
		}
	})
	if err != nil {
		log.Fatalf("Failed to setup cron handler: %v", err)
	}

	log.Println("Bot starting up.")
	c.Start()
	wg.Wait()
}

func notifier(t time.Time) bool {
	var err error

	// if we don't have a (valid) menu, try again to download it
	if currentMenu == nil {
		log.Println("Attempting to retrieve this week's menu.")

		currentMenu, err = fetchMenu(t)
		if err != nil {
			log.Printf("Failed to retrieve menu: %v", err)
			return false
		}

		log.Println("Downloaded the menu successfully.")
	}

	// we have a menu, report it to slack!
	log.Println("Reporting menu to Slack.")
	sendSlackNotification(t)

	return true
}

func fetchMenu(t time.Time) (*menu, error) {
	body, err := fetchMenuImage()
	if err != nil {
		return nil, fmt.Errorf("failed to determine menu image URL: %v", err)
	}
	defer body.Close()

	tmpFile := "/tmp/foo.png"
	defer os.Remove(tmpFile)

	err = prepareImage(body, tmpFile)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare image: %v", err)
	}

	text, err := ocr(tmpFile)
	if err != nil {
		return nil, fmt.Errorf("failed to OCR image: %v", err)
	}

	m, err := parseMenu(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse menu: %v", err)
	}

	if m.beginDate.IsZero() || m.endDate.IsZero() {
		return nil, errors.New("could not detect date range for the menu")
	}

	if t.Before(m.beginDate) || t.After(m.endDate) {
		return nil, fmt.Errorf("menu is out of date, valid for %s to %s", m.beginDate.Format("Jan 02"), m.endDate.Format("Jan 02"))
	}

	return m, nil
}

func formatItem(item menuItem) string {
	return fmt.Sprintf("*%s*\n%s _[%s]_", item.title, item.subtitle, item.price)
}

func sendSlackNotification(t time.Time) {
	weekdayIdx := weekdayIndex(t)

	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "M I", Value: formatItem(currentMenu.weekdays[weekdayIdx].M1)})
	attachment.AddField(slack.Field{Title: "M II", Value: formatItem(currentMenu.weekdays[weekdayIdx].M2)})

	payload := slack.Payload{
		Username:    fmt.Sprintf("Tagesmenü für %s", formatDate(t)),
		Channel:     slackChannel,
		IconEmoji:   slackUsername,
		Attachments: []slack.Attachment{attachment},
	}

	errs := slack.Send(slackURL, "", payload)
	if len(errs) > 0 {
		log.Printf("Failed to send Slack notification: %v", errs)
	}
}
