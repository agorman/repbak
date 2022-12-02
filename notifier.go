package repbak

// Notifier defines a notification method.
type Notifier interface {
	// Notify sends a notification
	Notify(error) error
}
