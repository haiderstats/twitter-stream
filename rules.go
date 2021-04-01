package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	models "twitter-streams/models"

	"github.com/joho/godotenv"
)

var client = &http.Client{}

func getRules(cli *http.Client, url_path string) (models.GetResponse, error) {

	req, _ := http.NewRequest("GET", url_path, nil)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))

	resp, err := cli.Do(req)

	if err != nil {
		log.Fatal("Error with request: ", err)
	}
	var rules models.GetResponse
	if resp.StatusCode == http.StatusOK {
		var rules models.GetResponse
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &rules)

		return rules, nil
	}
	return rules, fmt.Errorf("Error: %s", resp.Status)
}

func deleteRules(cli *http.Client, url_path string, ids models.IdList) error {
	toDelete, _ := json.Marshal(models.DeleteRules{Delete: ids})
	req, _ := http.NewRequest("POST", url_path, bytes.NewBuffer(toDelete))
	req.Header.Add("Authorization", "Bearer "+os.Getenv(("BEARER_TOKEN")))
	req.Header.Add("Content-type", "application/json")

	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusOK {
		var metaDelete models.DeleteResponse
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &metaDelete)

		if metaDelete.Errors != nil {
			for _, errors := range *metaDelete.Errors {
				for _, e := range errors.Errors {
					log.Println("Error in Delete:", e.Message)
				}
			}
		}
		return nil
	}
	return fmt.Errorf("Error: %s", resp.Status)
}

func createRules(cli *http.Client, url_path string, rules []models.AddRule) ([]models.Rule, error) {
	toAdd, _ := json.Marshal(models.CreateRules{Add: rules})

	req, _ := http.NewRequest("POST", url_path, bytes.NewBuffer(toAdd))
	req.Header.Add("Authorization", "Bearer "+os.Getenv(("BEARER_TOKEN")))
	req.Header.Add("Content-type", "application/json")

	resp, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var response models.CreateResponse

	switch resp.StatusCode {
	case 201:
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &response)
		if response.Errors != nil {
			for _, errors := range *response.Errors {
				log.Printf("Rule: %s failed due to %s", errors.Value, errors.Title)
			}
		}
		return response.Data, nil

	case 200:
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &response)
		//200 implies errors...unfortunately
		var errstrings []string
		for _, errors := range *response.Errors {
			var detail string
			for _, d := range errors.Details {
				detail += d + ". "
			}
			err = fmt.Errorf("The value '%s', failed with reason '%s': %s", errors.Value, errors.Title, detail)
			errstrings = append(errstrings, err.Error())
		}
		return nil, fmt.Errorf(strings.Join(errstrings, "\n"))

	}
	return nil, fmt.Errorf("Error : %s", resp.Status)
}

func startUp() []models.Rule {
	rules_url, err := url.Parse("https://api.twitter.com/2/tweets/search/stream/rules")
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load("tweet.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rules, err := getRules(client, rules_url.String())
	if err != nil {
		log.Fatal("Failed to retrieve rules. Error:", err)
	}

	var ids []string
	for _, rule := range rules.Data {
		ids = append(ids, rule.Id)
	}

	if len(ids) > 0 {
		err = deleteRules(client, rules_url.String(), models.IdList{Ids: ids})
		if err != nil {
			log.Fatal("Failed to delete Rules: ", err)
		}
	}

	rulesToAdd := make([]models.AddRule, 1)
	rulesToAdd[0] = models.AddRule{Value: `"trees" has:images lang:en -is:retweet`, Tag: "test tree images"}
	created, err := createRules(client, rules_url.String(), rulesToAdd)

	if err != nil {
		log.Fatal(err)
	}

	return created
}

func main() {
	_ = startUp()

	stream, err := url.Parse("https://api.twitter.com/2/tweets/search/stream")
	q := stream.Query()
	q.Set("tweet.fields", "attachments")
	q.Set("expansions", "attachments.media_keys")
	q.Set("media.fields", "url")
	stream.RawQuery = q.Encode()
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("GET", stream.String(), nil)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))

	resp, err := client.Do(req)

	dec := json.NewDecoder(resp.Body)

	for {
		var t models.Tweet
		dec.Decode(&t)
		log.Println(t)
	}

	return
}
