// Copyright 2018-2019 Arthur Bailão. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE.md file.

package main

import (
	"encoding/base64"
	"fmt"
	"os"

	kcl "github.com/arthurbailao/aws-kcl"
)

// WriteToFile implements aws-kcl RecordProcessor interface
type WriteToFile struct {
	outfile *os.File
}

// Initialize open the  output file
func (ec *WriteToFile) Initialize(shardID string) error {

	var err error
	ec.outfile, err = os.Create(fmt.Sprintf("/tmp/%s.txt", shardID))
	if err != nil {
		return err
	}

	fmt.Fprintf(ec.outfile, "init: %s\n", shardID)
	return nil
}

// ProcessRecords writes decoded data to file
func (ec *WriteToFile) ProcessRecords(records []*kcl.Record, checkpointer *kcl.Checkpointer) error {

	var sn string
	for i := range records {
		decoded, err := base64.StdEncoding.DecodeString(records[i].Data)
		if err != nil {
			return err
		}
		fmt.Fprintf(ec.outfile, "process: %s\n", decoded)
		sn = records[i].SequenceNumber
	}

	return checkpointer.Checkpoint(&sn)
}

// ShutdownRequested closes the file
func (ec *WriteToFile) ShutdownRequested(checkpointer *kcl.Checkpointer) error {
	fmt.Fprintln(ec.outfile, "shutdown")
	ec.outfile.Close()
	return checkpointer.Checkpoint(nil)
}

// Shutdown makes the last checkpoint
func (ec *WriteToFile) Shutdown(shutdownType kcl.ShutdownType, checkpointer *kcl.Checkpointer) error {
	if shutdownType == kcl.GracefulShutdown {
		if err := checkpointer.Checkpoint(nil); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var ec WriteToFile
	kcl.Run(&ec)
}
