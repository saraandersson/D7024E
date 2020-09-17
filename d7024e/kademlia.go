package d7024e

type Kademlia struct {
	network *Network
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func NewKademlia(network *Network) Kademlia{
	kademlia := Kademlia{}
	kademlia.network=network
	return kademlia
}
