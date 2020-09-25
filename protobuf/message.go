package protobuf

type Message struct {
	Text string `protobuf:"bytes,1,opt,name=Text" json:"Text,omitempty"`
}