package main

import (
	"fmt"

	"github.com/pquerna/ffjson/ffjson"
)

type PrivateRegistry struct{}

func (r *PrivateRegistry) GetAll() {
	b, err := get(fmt.Sprintf("%v://%s/%s/%s", protocol, srcRegistry, apiVer, searchUrl))
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
		r.GetTags(result.Results[i].Name)
	}
}

func (r *PrivateRegistry) GetTags(repositoryName string) {
	b, err := get(fmt.Sprintf("%s://%s/%s/repositories/%s/tags", protocol, srcRegistry, apiVer, repositoryName))
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
		//		pullCmd := fmt.Sprintf("docker pull %s/%s:%s", srcRegistry, repositoryName, k)
		//		fmt.Println(pullCmd)
		tagCmd := fmt.Sprintf("docker tag %s %s/%s:%s", value, dstRegistry, repositoryName, k)
		fmt.Println(tagCmd)
		pushCmd := fmt.Sprintf("docker push %s/%s:%s", dstRegistry, repositoryName, k)
		fmt.Println(pushCmd)
	}
}
