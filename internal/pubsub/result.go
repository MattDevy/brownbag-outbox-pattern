package pubsub

import "context"

// A PublishResult holds the result from a call to Publish.
type PublishResult struct {
	ready    chan struct{}
	serverID string
	err      error
}

// Ready returns a channel that is closed when the result is ready.
// When the Ready channel is closed, Get is guaranteed not to block.
func (r *PublishResult) Ready() <-chan struct{} { return r.ready }

// Get returns the server-generated message ID and/or error result of a Publish call.
// Get blocks until the Publish call completes or the context is done.
func (r *PublishResult) Get(ctx context.Context) (serverID string, err error) {
	// If the result is already ready, return it even if the context is done.
	select {
	case <-r.Ready():
		return r.serverID, r.err
	default:
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-r.Ready():
		return r.serverID, r.err
	}
}

// NewPublishResult creates a PublishResult.
func NewPublishResult() *PublishResult {
	return &PublishResult{ready: make(chan struct{})}
}

// SetPublishResult sets the server ID and error for a publish result and closes
// the Ready channel.
func SetPublishResult(r *PublishResult, sid string, err error) {
	r.serverID = sid
	r.err = err
	close(r.ready)
}
