package relay

type Servicer interface {
	Enable(relayID, itemID int) error
	Disable(relayID, itemID int) error
	Toggle(relayID, itemID int) error
}
