package reload

import "context"

// Notifier defines the interface for different notification mechanisms.
type Notifier interface {
	// WatchForChange monitors for changes and triggers the provided onReload function when a reload is needed.
	WatchForChange(ctx context.Context, onReload func(file string) error) error

	// Close stops any internal resources associated with the notifier.
	Close() error
}
