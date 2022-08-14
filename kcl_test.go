package kcl

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
)

type testProcessor struct {
}

func (t testProcessor) Initialize(string) error {
	return nil
}

func (t testProcessor) ProcessRecords(records []*Record, c *Checkpointer) error {
	sn := records[len(records)-1].SequenceNumber
	err := c.Checkpoint(&sn)
	return err
}

func (t testProcessor) ShutdownRequested(*Checkpointer) error {
	return nil
}

func (t testProcessor) LeaseLost() error {
	return nil
}

func (t testProcessor) ShardEnded(*Checkpointer) error {
	return nil
}

type printfMock struct {
	mock.Mock
}

func (m *printfMock) Printf(format string, a ...any) (int, error) {
	args := m.Called(format, a)
	return args.Int(0), args.Error(1)
}

func mockStdin(t *testing.T, input string) (funcDefer func(), err error) {
	t.Helper()

	oldOsStdin := os.Stdin
	tmpfile, err := ioutil.TempFile(t.TempDir(), t.Name())

	if err != nil {
		return nil, err
	}

	content := []byte(input)

	if _, err := tmpfile.Write(content); err != nil {
		return nil, err
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return nil, err
	}

	// Set stdin to the temp file
	os.Stdin = tmpfile

	return func() {
		// clean up
		os.Stdin = oldOsStdin
		os.Remove(tmpfile.Name())
	}, nil
}

func TestRunInitialize(t *testing.T) {
	var input = `{"action":"initialize","shardId":"shardId-000000000001"}` + "\n"
	var initializeStatus = `{"action":"status","responseFor":"initialize"}`

	f, err := mockStdin(t, input)
	if err != nil {
		t.Errorf("failed to mock os.Stdin %s", err)
	}

	defer f()
	m := printfMock{}
	printf = m.Printf

	m.On("Printf", "\n%s\n", []any{initializeStatus}).Return(0, nil)
	Run(testProcessor{})
	m.AssertExpectations(t)
}

func TestRunProcessRecords(t *testing.T) {
	var input = `{"action":"processRecords","records":[{"data":"bWVvdw==","partitionKey":"cat","sequenceNumber":"456"}]}` + "\n"
	var processRecordsStatus = `{"action":"status","responseFor":"processRecords"}`
	var checkpointAction = `{"action":"checkpoint","sequenceNumber":"456"}`

	f1, err := mockStdin(t, input)
	if err != nil {
		t.Errorf("failed to mock os.Stdin %s", err)
	}
	defer f1()

	m := printfMock{}
	printf = m.Printf

	m.On("Printf", "\n%s\n", []any{processRecordsStatus}).Return(0, nil)

	var f2 func()
	m.On("Printf", "\n%s\n", []any{checkpointAction}).Return(0, nil).Run(func(_ mock.Arguments) {
		f2, err = mockStdin(t, `{"action":"checkpoint","checkpoint":"456"}`)
		if err != nil {
			t.Errorf("failed to mock os.Stdin %s", err)
		}
	})

	Run(testProcessor{})
	defer f2()
	m.AssertExpectations(t)
}

func TestRunAllPhases(t *testing.T) {
	var initializeInput = `{"action":"initialize","shardId":"shardId-000000000001"}` + "\n"
	var initializeStatus = `{"action":"status","responseFor":"initialize"}`
	var processRecordsInput = `{"action":"processRecords","records":[{"data":"bWVvdw==","partitionKey":"cat","sequenceNumber":"456"}]}` + "\n"
	var processRecordsStatus = `{"action":"status","responseFor":"processRecords"}`
	var checkpointAction = `{"action":"checkpoint","sequenceNumber":"456"}`
	var shardEndedInput = `{"action":"shardEnded"}` + "\n"
	var shardEndedStatus = `{"action":"status","responseFor":"shardEnded"}`

	var funcs []func()
	f, err := mockStdin(t, initializeInput)
	if err != nil {
		t.Errorf("failed to mock os.Stdin %s", err)
	}
	funcs = append(funcs, f)

	m := printfMock{}
	printf = m.Printf

	m.On("Printf", "\n%s\n", []any{initializeStatus}).Return(0, nil).Run(func(args mock.Arguments) {
		f, err := mockStdin(t, processRecordsInput)
		if err != nil {
			t.Errorf("failed to mock os.Stdin %s", err)
		}
		funcs = append(funcs, f)
	})

	m.On("Printf", "\n%s\n", []any{processRecordsStatus}).Return(0, nil).Run(func(args mock.Arguments) {
		f, err := mockStdin(t, shardEndedInput)
		if err != nil {
			t.Errorf("failed to mock os.Stdin %s", err)
		}
		funcs = append(funcs, f)

	})
	m.On("Printf", "\n%s\n", []any{checkpointAction}).Return(0, nil).Run(func(_ mock.Arguments) {
		f, err = mockStdin(t, `{"action":"checkpoint","checkpoint":"456"}`)
		if err != nil {
			t.Errorf("failed to mock os.Stdin %s", err)
		}
		funcs = append(funcs, f)
	})

	m.On("Printf", "\n%s\n", []any{shardEndedStatus}).Return(0, nil)

	Run(testProcessor{})

	for _, f := range funcs {
		f()
	}
	m.AssertExpectations(t)
}
