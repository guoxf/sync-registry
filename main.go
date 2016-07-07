package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pquerna/ffjson/ffjson"
)

var (
	srcRegistry = "127.0.0.1:5000"
	dstRegistry = "reg.mydomain.com"
	apiVer      = "v1"
	userName    = "admin"
	pwd         = "Harbor12345"
)

const (
	tagUrl    = "repositories/%s/tags"
	searchUrl = "search"
)

type ImageInfo struct {
	Description string
	Name        string
}

type SearchResult struct {
	NumResults int `json:"num_results"`
	Query      string
	Results    []ImageInfo
}

func init() {
	flag.StringVar(&srcRegistry, "regsrc", "127.0.0.1:5000", "src registry host")
	flag.StringVar(&dstRegistry, "regdst", "reg.mydomain.com", "src registry host")
	flag.StringVar(&apiVer, "apiver", "v1", "api version")
	flag.Parse()
}

func main() {
	resp, err := http.Get(fmt.Sprintf("http://%s/%s/%s", srcRegistry, apiVer, searchUrl))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var result SearchResult
	err = ffjson.Unmarshal(b, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range result.Results {
		getTags(result.Results[i].Name)
	}
}

func getTags(repositoryName string) {
	resp, err := http.Get(fmt.Sprintf("http://%s/%s/repositories/%s/tags", srcRegistry, apiVer, repositoryName))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var tags = make(map[string]string)
	err = ffjson.Unmarshal(b, &tags)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, value := range tags {
		pullCmd := fmt.Sprintf("docker pull %s/%s:%s", srcRegistry, repositoryName, k)
		fmt.Println(pullCmd)
		tagCmd := fmt.Sprintf("docker tag %s %s/%s:%s", value, srcRegistry, repositoryName, k)
		fmt.Println(tagCmd)
		pushCmd := fmt.Sprintf("docker push %s/%s:%s", dstRegistry, repositoryName, k)
		fmt.Println(pushCmd)
	}
}
