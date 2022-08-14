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

	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	printf("\n%s\n", string(b))

	msg := readMessage(os.Stdin)
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
