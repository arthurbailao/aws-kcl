package kcl

// RecordProcessor ...
type RecordProcessor interface {
	Initialize(string) error

	ProcessRecords([]*Record, *Checkpointer) error

	ShutdownRequested(*Checkpointer) error

	Shutdown(ShutdownType, *Checkpointer) error
}

// Record ...
type Record struct {
	Data           string `json:"data"`
	PartitionKey   string `json:"partitionKey"`
	SequenceNumber string `json:"sequenceNumber"`
}
