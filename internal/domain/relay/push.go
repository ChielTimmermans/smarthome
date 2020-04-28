package relay

type Pusher interface {
	Enable(relayID, itemID int) error
	Disable(relayID, itemID int) error
}
