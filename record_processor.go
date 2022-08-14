package kcl

// RecordProcessor ...
type RecordProcessor interface {
	// Initialize is called once by the KCL before any calls to processRecords.
	// Any initialization logic for record processing can go here.
	Initialize(string) error

	// ProcessRecords is called by KCL with a list of records to be processed
	// and checkpointed. The checkpointer can optionally be used to checkpoint
	// a particular sequence number (from a record).
	ProcessRecords([]*Record, *Checkpointer) error

	// ShutdownRequested is called by KCL to indicate that this record processor
	// should shut down. After the shutdown requested operation is complete,
	// there will not be any more calls to any other functions of this record
	// processor. Clients should not attempt to checkpoint because the lease has
	// been lost by this Worker.
	ShutdownRequested(*Checkpointer) error

	// LeaseLost is called by the KCL to indicate that this record processor should shut down.
	// After the lease lost operation is complete, there will not be any more calls to
	// any other functions of this record processor. Clients should not attempt to
	// checkpoint because the lease has been lost by this Worker.
	LeaseLost() error

	// ShardEnded is called by the KCL to indicate that this record processor should shutdown.
	// After the shard ended operation is complete, there will not be any more calls to
	// any other functions of this record processor. Clients are required to checkpoint
	// at this time. This indicates that the current record processor has finished
	// processing and new record processors for the children will be created.
	ShardEnded(*Checkpointer) error
}

// Record ...
type Record struct {
	Data           string `json:"data"`
	PartitionKey   string `json:"partitionKey"`
	SequenceNumber string `json:"sequenceNumber"`
}
