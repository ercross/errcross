package errcross

//ErrcrossRepository provides a repository for Errcross

//Find finds a key in the repository, returns the Errcross instance associated with the provided key if found, error otherwise

//Store saves an Errcross instance into the repository
type ErrcrossRepository interface {
	Find(key string) (*Errcross, error)
	Store(e *Errcross) error
}
