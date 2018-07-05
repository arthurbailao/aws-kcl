package kcl

import (
	"fmt"
	"os"
)

// Run ...
func Run(processor RecordProcessor) {
	checkpointer := &Checkpointer{}

	for {
		msg := ReadMessage()

		if msg == nil {
			break
		}

		var err error
		switch {
		case msg.Action == "processRecords":
			err = processor.ProcessRecords(msg.Records, checkpointer)

		case msg.Action == "initialize":
			err = processor.Initialize(*msg.ShardID)

		case msg.Action == "shutdown":
			shutdownType := GracefulShutdown
			if msg.Reason == nil || *msg.Reason == "ZOMBIE" {
				shutdownType = ZombieShutdown
			}
			err = processor.Shutdown(shutdownType, checkpointer)

		case msg.Action == "shutdownRequested":
			err = processor.ShutdownRequested(checkpointer)

		default:
			err = fmt.Errorf("Unsupported KCL action: %s", msg.Action)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		WriteStatus(msg.Action)
	}
}
