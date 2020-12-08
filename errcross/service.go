package errcross

//ErcrossService is an interface that specifies the services oferred by Errcross

//Find finds the url associated with Key
//returns a pointer to the Errcross instance if the key is found in repository

//Store saves an Errcross instance into the repository, returns an error if save operation fails
type ErrcrossService interface {
	Find(key string) (*Errcross, error)
	Store(errcross *Errcross) error
}
