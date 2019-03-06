package kcl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// BrokenReader will return an error when Read is called:
type BrokenReader struct{}

func (r BrokenReader) Read(b []byte) (int, error) {
	return 0, fmt.Errorf("BrokenReader")
}

// EofReader will return an io.EOF error when Read is called:
type EofReader struct{}

func (r EofReader) Read(b []byte) (int, error) {
	return 0, io.EOF
}

func strP(str string) *string {
	return &str
}

func TestReadMessage(t *testing.T) {
	t.Run("Should parse JSON messages correctly", func(t *testing.T) {
		fakeMsg := message{
			Action:  "FakeAction",
			ShardID: strP("FakeShardID"),
			Records: []*Record{
				&Record{
					Data:           "FakeData",
					PartitionKey:   "FakePartitionKey",
					SequenceNumber: "FakeSequenceNumber",
				},
			},
			Reason: strP("FakeReason"),
			Error:  strP("FakeError"),
		}
		rawMsg, _ := json.Marshal(fakeMsg)

		tests := []struct {
			desc   string
			input  io.Reader
			output *message
		}{
			{
				desc:   "normal message",
				input:  bytes.NewReader(rawMsg),
				output: &fakeMsg,
			},
			{
				desc:   "message with EOF",
				input:  EofReader{},
				output: nil,
			},
		}

		for _, test := range tests {
			msg, err := readMessage(test.input)
			if err != nil {
				t.Errorf(
					"%s test: unexpected error from readMessage(): %s",
					test.desc,
					err.Error(),
				)
				return
			}

			if diff := cmp.Diff(test.output, msg); diff != "" {
				t.Errorf(
					"%s test: unexpected value on BuildEvents output: (-want +got)\n%s",
					test.desc,
					diff,
				)
			}
		}
	})

	t.Run("Should report errors correctly", func(t *testing.T) {
		_, err := readMessage(BrokenReader{})
		if err == nil {
			t.Errorf("expected error from readMessage() but got nothing")
		}

		_, err = readMessage(strings.NewReader("not a valid json"))
		if err == nil {
			t.Errorf("expected invalid json error from readMessage() but got nothing")
		}
	})
}
