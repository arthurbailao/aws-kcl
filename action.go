package kcl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type message struct {
	Action  string    `json:"action"`
	ShardID *string   `json:"shardId"`
	Records []*Record `json:"records"`
	Reason  *string   `json:"reason"`
	Error   *string   `json:"error"`
}

type status struct {
	Action      string `json:"action"`
	ResponseFor string `json:"responseFor"`
}

// ShutdownType ...
type ShutdownType int

const (
	unknownShutdown ShutdownType = iota

	// GracefulShutdown ...
	GracefulShutdown

	// ZombieShutdown ...
	ZombieShutdown
)

var printf = fmt.Printf

func readMessage(rd io.Reader) *message {
	bio := bufio.NewReader(rd)
	var buffer bytes.Buffer
	for {
		line, more, err := bio.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			panic("Unable to read line from stdin " + err.Error())
		}
		buffer.Write(line)
		if more == false {
			break
		}
	}

	var msg message
	err := json.Unmarshal(buffer.Bytes(), &msg)
	if err != nil {
		panic("Failed to unmarshal json message " + err.Error())
	}
	return &msg
}

func writeStatus(action string) {
	s := status{"status", action}
	b, err := json.Marshal(s)
	if err != nil {
		panic("Failed to marshal status as json " + err.Error())
	}

	printf("\n%s\n", string(b))
}
