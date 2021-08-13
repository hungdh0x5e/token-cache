package tokencache

import "golang.org/x/oauth2"

// inflight is used to wait on some in-flight request from multiple goroutines.
type inflight struct {
	doneCh chan struct{}

	token *oauth2.Token
	err   error
}

func newInflight() *inflight {
	return &inflight{doneCh: make(chan struct{})}
}

// wait returns a channel that multiple goroutines can receive on. Once it returns
// a value, the inflight request is done and result() can be inspected.
func (i *inflight) wait() <-chan struct{} {
	return i.doneCh
}

// done can only be called by a single goroutine. It records the result of the
// inflight request and signals other goroutines that the result is safe to
// inspect.
func (i *inflight) done(token *oauth2.Token, err error) {
	i.token = token
	i.err = err
	close(i.doneCh)
}

// result cannot be called until the wait() channel has returned a value.
func (i *inflight) result() (*oauth2.Token, error) {
	return i.token, i.err
}
