// Copyright 2018-2019 Arthur Bailão. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE.md file.

package kcl

import (
	"encoding/json"
	"fmt"
	"os"
)

// Checkpointer ...
type Checkpointer struct {
}

type checkpoint struct {
	Action         string  `json:"action"`
	SequenceNumber *string `json:"sequenceNumber"`
}

// Checkpoint ...
func (cp *Checkpointer) Checkpoint(sequenceNumber *string) error {
	c := checkpoint{"checkpoint", sequenceNumber}

	str, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n%s\n", str)

	msg := readMessage()
	if msg == nil {
		fmt.Fprintf(os.Stderr, "Received EOF rather than checkpoint ack\n")
		os.Exit(1)
	} else if msg.Action != "checkpoint" {
		fmt.Fprintf(os.Stderr, "Received invalid checkpoint ack: %s\n", msg.Action)
		os.Exit(1)
	} else if msg.Error != nil {
		return fmt.Errorf(*msg.Error)
	}

	return nil
}
