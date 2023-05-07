package client

import (
	"buptcs/entity"
	"bytes"
	"encoding/json"
	"io/ioutil" // TODO(goverclock): fix warning
	"log"
	"net/http"
)

type Client interface {
	RequestLogin(string, string) bool
	RequestRegister(string, string) bool
}

type clientImpl struct {
	// logged bool
	client http.Client
}

func New() Client {
	return &clientImpl{}
}

// this implementation may change
// but interface shouldn't change
func (ci *clientImpl) RequestLogin(username string, passwd string) bool {
	// forward request to server
	loginfo := entity.LogInfo{Username: username, Passwd: passwd}
	body, err := json.Marshal(loginfo)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/login", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")	
	resp, err := ci.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body_content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var status entity.LogStatus
	err = json.Unmarshal(body_content, &status)
	if err != nil {
		log.Fatal(err)
	}

	return status.Success
}

// TODO(goverclock): chang all "Log" to "Reg"
func (ci *clientImpl) RequestRegister(username string, passwd string) bool {
	// forward request to server
	loginfo := entity.LogInfo{Username: username, Passwd: passwd}
	body, err := json.Marshal(loginfo)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/register", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")	
	resp, err := ci.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body_content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var status entity.LogStatus
	err = json.Unmarshal(body_content, &status)
	if err != nil {
		log.Fatal(err)
	}

	return status.Success
}
