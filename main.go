package main

import "flag"

var (
	srcRegistry = "127.0.0.1:5000"
	dstRegistry = "reg.mydomain.com"
	apiVer      = "v1"
	userName    = "admin"
	pwd         = "Harbor12345"
	orgname     = ""
	protocol    = "https"
)

const (
	tagUrl          = "repositories/%s/tags"
	repositoriesUrl = "repositories/%s"
	searchUrl       = "search"
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
	flag.StringVar(&protocol, "protocol", "https", "protocol")
	flag.StringVar(&srcRegistry, "regsrc", "hub.docker.com", "src registry host")
	flag.StringVar(&dstRegistry, "regdst", "reg.mydomain.com", "src registry host")
	flag.StringVar(&apiVer, "apiver", "v2", "api version")
	flag.StringVar(&orgname, "orgname", "rancher", "orgname")
	flag.Parse()
}

func main() {
	r := &DockerRegistry{
		results: make([]string, 0),
		m:       make(map[string][]string),
	}
	r.GetRepositories("")
	r.Save("rancher.sh")
}
