package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var gitlabPAToken string
var git *gitlab.Client

type Config struct {
	Feeds    []Feed
	Interval int
}

type Feed struct {
	FeedURL         string `yaml:"feed_url"`
	Name            string
	GitlabProjectID int `yaml:"gitlab_project_id"`
	Labels          []string
}

func (feed Feed) checkFeed(lastRun time.Time) {
	fp := gofeed.NewParser()
	rss, err := fp.ParseURL(feed.FeedURL)

	if err != nil {
		fmt.Printf("Unable to parse feed %s: \n %s", feed.Name, err)
		return
	}

	var newArticle []*gofeed.Item
	var oldArticle []*gofeed.Item
	for _, item := range rss.Items {
		var time *time.Time

		// Prefer updated time to published
		if item.UpdatedParsed != nil {
			time = item.UpdatedParsed
		} else {
			time = item.PublishedParsed
		}

		if time.After(lastRun) {
			newArticle = append(newArticle, item)
		} else {
			oldArticle = append(oldArticle, item)
		}
	}

	fmt.Printf("Feed Name: %s\n", feed.Name)
	fmt.Printf("Old Items: %d\n", len(oldArticle))
	fmt.Printf("New Items: %d\n", len(newArticle))

	for _, item := range newArticle {

		// Prefer description over content
		var body string
		if item.Description != "" {
			body = item.Description
		} else {
			body = item.Content
		}

		issueOptions := &gitlab.CreateIssueOptions{
			Title:       gitlab.String(item.Title),
			Description: gitlab.String(body),
			Labels:      feed.Labels,
		}
		_, _, err := git.Issues.CreateIssue(feed.GitlabProjectID, issueOptions)
		if err != nil {
			fmt.Printf("Unable to create Gitlab issue for %s \n %s", feed.Name, err)
		} else {
			fmt.Printf("Creating Gitlab Issue '%s' in project: %d'", issueOptions.Title, feed.GitlabProjectID)
		}
	}
}

func readConfig(path string) *Config {
	config := &Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Fatalln("Unable to parse config YAML \n ")
	}
	return config
}

func main() {
	var lastRun = time.Now()
	readEnv()
	git = gitlab.NewClient(nil, gitlabPAToken)

	config := readConfig("config.yaml")

	for {
		fmt.Printf("Running checks at %s\n", time.Now().Format(time.RFC850))
		for _, configEntry := range config.Feeds {
			configEntry.checkFeed(lastRun)
		}
		lastRun = time.Now()
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}

}

func readEnv() {
	if envGitlabAPIToken := os.Getenv("GITLAB_API_TOKEN"); envGitlabAPIToken == "" {
		panic("Could not find GITLAB_API_TOKEN specified as an environment variable")
	} else {
		gitlabPAToken = envGitlabAPIToken
	}
}
