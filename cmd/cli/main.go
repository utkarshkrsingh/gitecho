package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var helpMessage string = `Usage:
    gitecho <username>
Example:
    gitecho utkarshkrsingh`

type githubEvent struct {
    Id   string `json:"id"`
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Commits []struct {
			Author struct {
				Name string `json:"name"`
			} `json:"author" `
		} `json:"commits"`
	} `json:"payload"`
    CreatedAt string `json:"created_at"`
}

func main() {

	if len(os.Args) < 2 || len(os.Args) > 2 {
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	var username string = os.Args[1]
	if _, err := strconv.Atoi(username); err == nil {
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to fetch events: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var events []githubEvent
    if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
        log.Printf("failed to decode JSON response: %v", err)
        os.Exit(1)
    }

    for _, event := range events {
        
        parsedTime, err := time.Parse(time.RFC3339, event.CreatedAt)
        if err != nil {
            log.Printf("error parsing the time for event ID %s: %v\n",event.Id, err)
            continue
        }
        readableTime := parsedTime.Format("02 Jan 2006 15:04:05")

        switch event.Type {
        case "PushEvent":
            commitCount := len(event.Payload.Commits)
            if commitCount > 0 {
                fmt.Printf("- %s Pushed %d commits to %s\n",readableTime, commitCount, event.Repo.Name)
            }
        case "IssuesEvent":
            fmt.Printf("- %s Opened a new issue in %s\n", readableTime, event.Repo.Name)
        case "WatchEvent":
            fmt.Printf("- %s Starred %s\n", readableTime, event.Repo.Name)
        }
    }
}
