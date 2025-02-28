package main

import (
	"fmt"
	"log"
	"net/http"
)

type updateValidator interface {
	checkUpdateValidity() bool
}
type createValidator interface {
	checkCreateValidity() bool
}

var validBreeds map[string]string

func initValidBreedsMap() {
	validBreeds = make(map[string]string)
	resp, err := http.Get("https://api.thecatapi.com/v1/breeds")
	if err != nil {
		log.Fatal("error making request:", err)
	}
	defer resp.Body.Close()

	var breedsData []struct {
		Name string `json:"name"`
	}
	err = decodeBody(resp.Body, &breedsData)
	if err != nil {
		log.Fatal("error decoding response:", err)
	}

	for _, breed := range breedsData {
		validBreeds[breed.Name] = breed.Name
	}
}

func (m Mission) checkCreateValidity() error {
	if m.SpyID == 0 {
		return fmt.Errorf("Mission object should include reference to Spy ID")
	}
	if len(m.Targets) > 3 {
		return fmt.Errorf("Target amount is %d, expected 0-3", len(m.Targets))
	}
	return nil
}

func (spy Spy) checkCreateValidity() bool {
	if _, ok := validBreeds[spy.Breed]; ok {
		return true
	}
	return false
}

func (spy Spy) checkUpdateValidity() bool {
	if spy.Salary == 0 {
		return false
	}
	return true
}

func (m Mission) checkUpdateValidity() bool {
	if !m.CompleteState && m.SpyID <= 0 {
		return false
	} else if m.CompleteState == true {
		return true
	}
	return true
}
