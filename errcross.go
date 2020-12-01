package errcross

import (
	"errors"
	"github.com/bradialabs/shortid"
	errs "github.com/pkg/errors"
	"gopkg.in/dealancer/validate.v2"
	"time"
)

var (
	ErrKeyNotFound  = errors.New("key does not exist")
	ErrMalformedUrl = errors.New("url malformed")
)

type Errcross struct {

	//Key is short code associated with a URL
	Key string `json:"key" bson:"key"`

	//URL is the exact URL supplied by user
	URL string `json:"url" bson:"url" validate:"empty=false & format=url"`

	//Timestamp denotes the generation time (since epoch) associated with this instance
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
}

//errcrossService implements the ErrcrossService interface
//errcrossRepository enables to fetch data from database and save into database
type errcrossService struct {
	errcrossRepository ErrcrossRepository
}

func NewErrcrossService(repository ErrcrossRepository) ErrcrossService {
	return &errcrossService{
		repository,
	}
}

func (e *errcrossService) Find(key string) (*Errcross, error) {
	return e.errcrossRepository.Find(key)
}

func (e *errcrossService) Store(url string) error {
	errcross := Errcross{
		URL: url,
	}

	//validate that the url is not empty and is in a url format
	if err := validate.Validate(errcross); err != nil {
		return errs.Wrap(ErrMalformedUrl, "@ErrcrossService.Store")
	}
	errcross.Key = shortid.New().Generate()
	errcross.Timestamp = time.Now().UTC().Unix()
	return e.errcrossRepository.Store(&errcross)
}
