package main

import (
	"fmt"
	"net/http"
)

func checkHealth(url string) (ok bool, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	ok = res.StatusCode == 200
	return
}

func checkProjectId(token string, id uint64) (valid bool, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", *serviceUrlProjects, id), nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	valid = res.StatusCode == 200
	return
}

func checkSprintId(token string, id uint64) (valid bool, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", *serviceUrlSprints, id), nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	valid = res.StatusCode == 200
	return
}
