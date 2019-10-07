package gossiper

//SimpleMessage contains all relevant information about the message.
type SimpleMessage struct {
	OriginalName  string
	RelayPeerAddr string
	Contents      string
}

//RumorMessage contains all relevant information about the rumor message.
type RumorMessage struct {
	Origin string
	ID     uint32
	Text   string
}

//PeerStatus contains information about the status of the peer.
type PeerStatus struct {
	Identifier string
	NextID     uint32
}

//StatusPacket is the status.
type StatusPacket struct {
	Want []PeerStatus
}

//GossipPacket To provide compatibility with future versions
type GossipPacket struct {
	Simple *SimpleMessage
	Rumor  *RumorMessage
	Status *StatusPacket
}
