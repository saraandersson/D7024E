package d7024e
import{
	"fmt"
}
type Kademlia struct {
	network *Network
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
	contacts := kademlia.network.routingTable.FindClosestContacts(target.ID, 2)
	fmt.Println(contacts)
	/*for i=0;i<len(contacts);i++{

	}*/

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
