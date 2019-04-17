package kcl

import (
	"fmt"
	"os"
)

// Run ...
func Run(processor RecordProcessor) {
	checkpointer := &Checkpointer{}

	for {
		var err error

		msg, err := readMessage(os.Stdin)
		if err != nil {
			panic(err)
		}
		if msg == nil {
			break
		}

		switch msg.Action {
		case "processRecords":
			err = processor.ProcessRecords(msg.Records, checkpointer)

		case "initialize":
			err = processor.Initialize(*msg.ShardID)

		case "shutdown":
			shutdownType := GracefulShutdown
			if msg.Reason == nil || *msg.Reason == "ZOMBIE" {
				shutdownType = ZombieShutdown
			}
			err = processor.Shutdown(shutdownType, checkpointer)

		case "shutdownRequested":
			err = processor.ShutdownRequested(checkpointer)

		default:
			err = fmt.Errorf("Unsupported KCL action: %s", msg.Action)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		writeStatus(msg.Action)
	}
}
