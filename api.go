package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Actor struct {
	Login        string `json:"login"`
	DisplayLogin string `json:"display_login"`
	GravatarId   string `json:"gravatar_id"`
	Url          string `json:"url`
	AvatarUrl    string `json:"avatar_url"`
}

type Repo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Author struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Commits struct {
	Sha      string `json:"sha"`
	Author   Author `json:"author"`
	Message  string `json:"message"`
	Distinct bool   `json:"distinct"`
	Url      string `json:"url"`
}
type Payload struct {
	PushId       int    `json:"push_id"`
	Size         int    `json:"size"`
	DistinctSize int    `json:"distinct_size"`
	Ref          string `json:"ref"`
	Head         string `json:"head"`
	Before       string `json:"before"`
	Commits      []Commits
}
type Org struct {
	Id         int    `json:"id"`
	Login      string `json:"login"`
	GravatarId string `json:"gravatar_id"`
	Url        string `json:"url"`
	AvatarUrl  string `json:"avatar_url"`
}

//Eventcoding json
type Event struct {
	Id           string  `json:"id"`
	Type         string  `json:"type"`
	Actor        Actor   `json:"actor"`
	Repo         Repo    `json:"repo"`
	Payload      Payload `json:"payload"`
	Public       bool    `json:"public"`
	CreatedAt    string  `json:"created_at"`
	Organization Org     `json:"org"`
}

type Activity struct {
	Id        string `json:"id"`
	Person    string `json:"person"`
	Repo      string `json:"repo"`
	CreatedAt string `json:"created_at"`
}

func main() {
	e, err := getEvents("kinvolk")
	if err != nil {
		fmt.Printf("problem getting events: %v", err)
		return
	}

	err = writeActivity("event.json", e)
	if err != nil {
		fmt.Printf("problem writing the activities to file: %v", err)
		return
	}
}

func getEvents(org string) ([]Event, error) {
	// get the data
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/events", org))
	if err != nil {
		return nil, fmt.Errorf("error getting the data from github: %v", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	var e []Event
	//saving the data to the new format
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		return nil, fmt.Errorf("problem unmarshalling data: %v", err)
	}
	return e, nil
}

func writeActivity(outfile string, e []Event) error {
	var activities []Activity
	for _, a := range e {
		act := Activity{
			Id:        a.Id,
			Person:    a.Actor.Login,
			Repo:      a.Repo.Name,
			CreatedAt: a.CreatedAt,
		}
		activities = append(activities, act)
	}

	b, err := json.Marshal(activities)
	if err = ioutil.WriteFile(outfile, b, 0666); err != nil {
		return fmt.Errorf("failed to write the file in json format: %v", err)
	}

	return nil
}
