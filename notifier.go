package repbak

// Notifier defines a notification method.
type Notifier interface {
	// Notify sends a notification
	Notify(stat Stat) error

	// Notify History sends a notification with the backup history.
	NotifyHistory(map[string][]Stat) error
}
