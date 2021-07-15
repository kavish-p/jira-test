package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type TransitionDetail struct {
	TransitionDateCreated string `json:"transitionDate"`
	FromString            string `json:"fromString"`
	ToString              string `json:"toString"`
	AuthorDisplayName     string `json:"authorDisplayName"`
}

type JIRAIssueSummary struct {
	Key             string `json:"key"`
	JIRADateCreated string `json:"dateCreated"`
	Transitions     []struct {
		TransitionDateCreated string `json:"transitionDate"`
		FromString            string `json:"fromString"`
		ToString              string `json:"toString"`
		AuthorDisplayName     string `json:"authorDisplayName"`
	} `json:"transitions"`
}

type JIRAIssue struct {
	Key    string `json:"key"`
	Fields struct {
		JIRADateCreated string `json:"created"`
	} `json:"fields"`
	Changelog struct {
		Histories []struct {
			Author struct {
				DisplayName string `json:"displayName"`
			} `json:"author"`
			TransitionDateCreated string `json:"created"`
			Items                 []struct {
				Field      string `json:"field"`
				FromString string `json:"fromString"`
				ToString   string `json:"toString"`
			} `json:"items"`
		} `json:"histories"`
	} `json:"changelog"`
}

func main() {
	// Get("https://track.appspace.com/rest/api/2/serverInfo")

	resp := GetIssue("AP-20504")

	var issue JIRAIssue

	err := json.Unmarshal([]byte(resp), &issue)
	if err != nil {
		panic(err)
	}

	jiraIssueSummary := JIRAIssueSummary{
		Key:             issue.Key,
		JIRADateCreated: issue.Fields.JIRADateCreated,
	}

	for _, history := range issue.Changelog.Histories {
		for _, transition := range history.Items {
			if transition.Field == "status" {

				jiraIssueSummary.Transitions = append(jiraIssueSummary.Transitions, TransitionDetail{
					TransitionDateCreated: history.TransitionDateCreated,
					FromString:            transition.FromString,
					ToString:              transition.ToString,
					AuthorDisplayName:     history.Author.DisplayName,
				})
			}
		}
	}

	fmt.Println(jiraIssueSummary)

	fmt.Println("")
	ProcessTransitions(jiraIssueSummary)

}

func GetIssue(issueID string) string {
	url := "https://track.appspace.com/rest/api/2/issue/" + issueID + "?expand=changelog"
	bearer := "Bearer " + "MjQ5OTQxMDExNTg0OhrDOXXqJrnFXx/LdURjKjOtZWnc"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	// log.Println(string([]byte(body)))
	return string([]byte(body))
}

func ProcessTransitions(summary JIRAIssueSummary) {

	fmt.Println(summary.Key)

	for _, transition := range summary.Transitions {
		fmt.Println(transition)
	}

}
