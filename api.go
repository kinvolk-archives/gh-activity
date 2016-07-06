package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// get the data
	resp, err := http.Get("https://api.github.com/users/kinvolk/events")
	if err != nil {
		fmt.Printf("error getting the data from github", err)
		return
	}

	type Actor struct {
		Id           int    `json:"id"`
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

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	var e []Event
	//saving the data to the new format
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		fmt.Printf("problem unmarshalling data: %v\n", err)
		return
	}

	fmt.Printf("marshalled data:\n %#v\n", e)

	type Activity struct {
		Id        string `json:"id"`
		Person    string `json:"person"`
		Repo      string `json:"repo"`
		CreatedAt string `json:"created_at"`
	}

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

	if err = ioutil.WriteFile("myjson.json", b, 0666); err != nil {
		fmt.Printf("Failed to write the file in json format", err)
	}
}
