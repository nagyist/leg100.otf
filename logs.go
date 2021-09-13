package ots

import "fmt"

type Logs []byte

type GetLogOptions struct {
	// The maximum number of bytes of logs to return to the client
	Limit int `schema:"limit"`

	// The start position in the logs from which to send to the client
	Offset int `schema:"offset"`
}

func (l Logs) Get(opts GetLogOptions) ([]byte, error) {
	if len(l) == 0 {
		return nil, nil
	}

	if opts.Offset > len(l) {
		return nil, fmt.Errorf("offset cannot be bigger than total logs")
	}

	if opts.Limit > MaxPlanLogsLimit {
		opts.Limit = MaxPlanLogsLimit
	}

	// Ensure specified chunk does not exceed slice length
	if (opts.Offset + opts.Limit) > len(l) {
		opts.Limit = len(l) - opts.Offset
	}

	return l[opts.Offset:(opts.Offset + opts.Limit)], nil
}

func (l *Logs) Append(logs []byte, opts UploadLogsOpts) {
	if opts.Start {
		// Add start marker
		*l = []byte{byte(2)}
	}

	*l = append(*l, logs...)

	if opts.End {
		// Add end marker
		*l = append(*l, byte(3))
	}
}
