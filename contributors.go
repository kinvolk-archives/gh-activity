package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Contributor struct {
	Login          string `json:"login"`
	Id             int    `json:"id"`
	AvatarUrl      string `json:"avatar_url"`
	GravatarId     string `json:"gravatar_id"`
	Url            string `json:"url"`
	HtmlUrl        string `json:"html_url"`
	Followers      string `json:"followers_url"`
	Following      string `json:"following_url"`
	Gists          string `json:"gists_url"`
	Starred        string `json:"starred_url"`
	Subscriptions  string `json:"subscriptions_url"`
	Organizations  string `json:"organizations_url"`
	Repos          string `json:"respos_url"`
	Events         string `json:"events_url"`
	ReceivedEvents string `json:"received_events_url"`
	Type           string `json:"type"`
	SiteAdmin      bool   `json:"site_admin"`
	Contributions  int    `json:"Contributions"`
}

var (
	org        string
	repo       string
	outputfile string

	Cmd = &cobra.Command{
		Use:   "gh-activity --org=[organization] --repo=[repository]",
		Short: "Get contributors in a repository",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if org == "" {
				fmt.Println("organization is not set")
				cmd.Help()
				return
			}
			if repo == "" {
				fmt.Println("repository not set")
				cmd.Help()
				return
			}
			e, err := getContributors(org, repo)
			if err != nil {
				fmt.Printf("Porblem getting the contributors: %v", err)
				return
			}

			err = writeContributors(outputfile, e)
			if err != nil {
				fmt.Printf("Problem writing the contributors to file: %v", err)
				return
			}

		},
	}
)

func init() {
	Cmd.PersistentFlags().StringVar(&org, "org", "", "Github organization")
	Cmd.PersistentFlags().StringVar(&repo, "repo", "", "Github repository")
	Cmd.PersistentFlags().StringVar(&outputfile, "out", "Contributions.json", "file to saves contributors")
}

func main() {
	if err := Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func getContributors(org string, repo string) ([]Contributor, error) {
	//get the data
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", org, repo))
	if err != nil {
		return nil, fmt.Errorf("Error getting the data from github: %v", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	var e []Contributor
	//saving the data to the new format
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		return nil, fmt.Errorf("problem unmarshalling data: %v", err)
	}
	return e, nil
}

func writeContributors(outputfile string, e []Contributor) error {
	var contributors []Contributor
	for _, a := range e {
		cntr := Contributor{
			Login:          a.Login,
			Id:             a.Id,
			AvatarUrl:      a.AvatarUrl,
			GravatarId:     a.GravatarId,
			Url:            a.Url,
			HtmlUrl:        a.HtmlUrl,
			Followers:      a.Followers,
			Following:      a.Following,
			Gists:          a.Gists,
			Starred:        a.Starred,
			Subscriptions:  a.Subscriptions,
			Organizations:  a.Organizations,
			Repos:          a.Repos,
			Events:         a.Events,
			ReceivedEvents: a.ReceivedEvents,
			Type:           a.Type,
			SiteAdmin:      a.SiteAdmin,
			Contributions:  a.Contributions,
		}
		contributors = append(contributors, cntr)
	}

	b, err := json.Marshal(contributors)
	if err = ioutil.WriteFile(outputfile, b, 0666); err != nil {
		return fmt.Errorf("faild to write the file in json format: %v", err)
	}

	return nil
}
