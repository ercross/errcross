package errcross

//ErrcrossSerializer is an interface that contains method to serialize and deserialize Errcross instance

//Decode takes a []byte, unmarshalls the slice and saves the resulting data into the empty receiver type
//Encode encodes the data contained in the receiver type, marshalls it and returns the marshalled data as a []byte
type ErrcrossSerializer interface {
	Decode(input []byte) (*Errcross, error)
	Encode(input *Errcross) ([]byte, error)
}
