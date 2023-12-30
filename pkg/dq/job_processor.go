package dq

type JobProcessorFunc func(*JobDecoder)

type JobDecoder struct {
	Name   string
	Key    int64
	Value  string
	Commit func(decoder *JobDecoder)
}
