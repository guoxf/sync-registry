package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"strings"

	"github.com/pquerna/ffjson/ffjson"
)

type DockerRegistry struct {
	results []string
	m       map[string][]string
}

func (r *DockerRegistry) genGetRepositoriesUrl(url string) string {
	if url != "" {
		return url
	}

	return fmt.Sprintf("%s://", protocol) + path.Join(srcRegistry, apiVer, "repositories", orgname)
}

func (r *DockerRegistry) genGetTagsUrl(repositoryName string) string {
	return fmt.Sprintf("%s://", protocol) + path.Join(srcRegistry, apiVer, "repositories", orgname, repositoryName, "tags")
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *DockerRegistry) Get(url string) (*Response, error) {
	b, err := get(r.genGetRepositoriesUrl(url))
	if err != nil {
		return nil, err
	}
	var response Response
	err = ffjson.Unmarshal(b, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type Response struct {
	Count    int
	Next     string
	Previous string
	Results  []map[string]interface{}
}

func (r *DockerRegistry) GetRepositories(url string) {
	repositories, err := r.Get(r.genGetRepositoriesUrl(url))
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, m := range repositories.Results {
		r.GetTags(m["name"].(string))
	}
	if repositories.Next != "" {
		r.GetRepositories(repositories.Next)
	}
}

func (r *DockerRegistry) GetTags(repositoryName string) {
	if _, ok := r.m[repositoryName]; !ok {
		r.m[repositoryName] = make([]string, 0)
	}
	tags, err := r.Get(r.genGetTagsUrl(repositoryName))
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, m := range tags.Results {
		name := m["name"].(string)
		if strings.Contains(name, "-rc") {
			fmt.Println(name)
			continue
		}
		r.m[repositoryName] = append(r.m[repositoryName], m["name"].(string))
		src := fmt.Sprintf("%v/%v:%v", orgname, repositoryName, m["name"])
		dst := fmt.Sprintf("%v/%v/%v:%v", dstRegistry, orgname, repositoryName, m["name"])

		pullCmd := fmt.Sprintf("docker pull %s/%s", srcRegistry, src)
		r.results = append(r.results, pullCmd)
		fmt.Println(pullCmd)

		tagCmd := fmt.Sprintf("docker tag %v %v", src, dst)
		r.results = append(r.results, tagCmd)
		fmt.Println(tagCmd)

		pushCmd := fmt.Sprintf("docker push %v", dst)
		r.results = append(r.results, pushCmd)
		return
	}
}

func (r *DockerRegistry) SaveBeautiful(fileName string) {
	imageNames := make([]string, 0)
	count := 0
	for k, v := range r.m {
		if len(v) == 0 {
			continue
		}
		imageNames = append(imageNames, fmt.Sprintf("%v:%v", k, v[0]))
		count++
		if count%5 == 0 {
			imageNames = append(imageNames, "\\ \n ")
		}
	}
	tpl := `
function pullRancherImage(){
	images=(%v)
	for imageName in ${images[@]} ; do
		docker pull {{orgname}}/$imageName
        docker tag {{orgname}}/$imageName {{regdst}}/{{orgname}}/$imageName
        docker push {{regdst}}/{{orgname}}/$imageName
		docker rmi {{orgname}}/$imageName
		docker rmi {{regdst}}/{{orgname}}/$imageName
    done
}
	`
	tpl = fmt.Sprintf(tpl, strings.Join(imageNames, " "))
	tpl = strings.Replace(tpl, "{{orgname}}", orgname, -1)
	tpl = strings.Replace(tpl, "{{regdst}}", dstRegistry, -1)
	ioutil.WriteFile(fileName, []byte(tpl), os.ModePerm)
}

func (r *DockerRegistry) Save(fileName string) {
	ioutil.WriteFile(fileName, []byte(strings.Join(r.results, "\n")), os.ModePerm)
}
