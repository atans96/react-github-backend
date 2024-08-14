package service

import (
	"backend/src/service/linguist"
	"backend/src/types"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/hashicorp/go-retryablehttp"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (f *Fetch) FetchFileName(repo types.Repository, userInfo UserInfo) {
	_, err := f.FetchRateLimit(userInfo.Token)
	if err != nil {
		panic(err)
	}
	makeNew, err := http.NewRequest("GET", "https://api.github.com/repos/"+userInfo.FullName+"/git/trees/master?recursive=1", nil)
	if err != nil {
		panic(err)
	}
	makeNew.Header.Set("Authorization", "Bearer "+userInfo.Token)
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.HTTPClient.Timeout = 50 * time.Second
	standardClient := retryClient.StandardClient() // *http.GQLClient
	response, err := standardClient.Do(makeNew)
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	go response.Body.Close()
	var fileName types.FileNameResponse
	err = json.Unmarshal(result, &fileName)
	if err != nil {
		panic(err)
	}
	trees := fileName.Tree
	var output types.RepoInfo
	var url string
	if len(repo.DefaultBranchRef.Name) == 0 {
		makeNew, err := http.NewRequest("GET", "https://api.github.com/repos/"+userInfo.FullName, nil)
		if err != nil {
			panic(err)
		}
		makeNew.Header.Set("Authorization", "Bearer "+userInfo.Token)
		makeNew.Header.Set("accept", "application/vnd.github.v3+json,application/vnd.github.mercy-preview+json,application/vnd.github.nebula-preview+json")
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = 10
		retryClient.HTTPClient.Timeout = 50 * time.Second
		standardClient := retryClient.StandardClient() // *http.GQLClient
		response, err := standardClient.Do(makeNew)
		result, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		go response.Body.Close()
		err = json.Unmarshal(result, &output)
		if err != nil {
			panic(err)
		}
		url = "https://raw.githubusercontent.com/" + userInfo.FullName + "/" + output.DefaultBranch + "/"
	} else {
		url = "https://raw.githubusercontent.com/" + userInfo.FullName + "/" + repo.DefaultBranchRef.Name + "/"
	}
	for _, tree := range trees {
		url = url + tree.Path
		makeNew, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
		makeNew.Header.Set("accept", "application/vnd.github.VERSION.raw")
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = 10
		retryClient.HTTPClient.Timeout = 50 * time.Second
		standardClient := retryClient.StandardClient() // *http.GQLClient
		response, err := standardClient.Do(makeNew)
		if err != nil {
			panic(err)
		}
		if !(response.StatusCode == 404) {
			result, err := ioutil.ReadAll(response.Body)
			if err != nil {
				panic(err)
			}
			_, allowed := linguist.LanguageType(string(result))
			if !(linguist.IsBinary(result) && linguist.IsVendored(string(result)) && linguist.IsDocumentation(string(result)) && allowed) {
				req := esapi.IndexRequest{
					Index:      os.Getenv("ES_INDEX_NAME"),
					DocumentID: strconv.Itoa(1),
					Body:       strings.NewReader(string(result)),
					Refresh:    "true",
				}
				res, err := req.Do(context.TODO(), ESClient)
				if err != nil {
					fmt.Printf("IndexRequest ERROR: %s\n", err)
				}
				if res.IsError() {
					fmt.Printf("%s ERROR indexing document ID=%d\n", res.Status(), 1)
				} else {
					// Deserialize the response into a map.
					var resMap map[string]interface{}
					if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
						fmt.Printf("Error parsing the response body: %s\n", err)
					} else {
						fmt.Println("\nIndexRequest() RESPONSE:")
						// Print the response status and indexed document version.
						fmt.Println("Status:", res.Status())
						fmt.Println("Result:", resMap["result"])
						fmt.Println("Version:", int(resMap["_version"].(float64)))
						fmt.Println("resMap:", resMap)
					}
				}
				err = res.Body.Close()
				if err != nil {
					panic(err)
				}
			}
		}
		err = response.Body.Close()
		if err != nil {
			panic(err)
		}
	}
}
