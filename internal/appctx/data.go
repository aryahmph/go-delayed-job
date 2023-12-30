package appctx

type WatcherData struct {
	Name   string
	Key    int64
	Value  string
	Commit func()
}
