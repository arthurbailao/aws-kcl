package main

import (
	"encoding/base64"
	"fmt"
	"os"

	kcl "github.com/arthurbailao/aws-kcl"
)

// WriteToFile implements aws-kcl RecordProcessor interface
type WriteToFile struct {
	outfile               *os.File
	currentSequenceNumber *string
}

// Initialize open the  output file
func (ec *WriteToFile) Initialize(shardID string) error {

	var err error
	ec.outfile, err = os.Create(fmt.Sprintf("/tmp/%s.txt", shardID))
	if err != nil {
		return err
	}

	fmt.Fprintln(ec.outfile, `{"action": "initialize"}`)

	return nil
}

// ProcessRecords writes decoded data to file
func (ec *WriteToFile) ProcessRecords(records []*kcl.Record, checkpointer *kcl.Checkpointer) error {
	for i := range records {
		decoded, err := base64.StdEncoding.DecodeString(records[i].Data)
		if err != nil {
			return err
		}
		fmt.Fprintln(ec.outfile, string(decoded))
		ec.currentSequenceNumber = &records[i].SequenceNumber
	}

	return checkpointer.Checkpoint(ec.currentSequenceNumber)
}

// ShutdownRequested closes the file
func (ec *WriteToFile) ShutdownRequested(checkpointer *kcl.Checkpointer) error {
	fmt.Fprintln(ec.outfile, `{"action": "shutdownRequested"}`)
	ec.outfile.Close()
	return checkpointer.Checkpoint(ec.currentSequenceNumber)
}

func (ec *WriteToFile) LeaseLost() error {
	return nil
}

// Shutdown makes the last checkpoint
func (ec *WriteToFile) ShardEnded(checkpointer *kcl.Checkpointer) error {
	fmt.Fprintln(ec.outfile, `{"action": "shardEnded"}`)
	ec.outfile.Close()
	return checkpointer.Checkpoint(ec.currentSequenceNumber)
}

func main() {
	var ec WriteToFile
	kcl.Run(&ec)
}
