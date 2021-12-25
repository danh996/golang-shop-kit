package internal

import (
	esv7 "github.com/elastic/go-elasticsearch/v7"
)

// NewElasticSearch instantiates the ElasticSearch client using configuration defined in environment variables.
func NewElasticSearch(conf esv7.Config) (es *esv7.Client, err error) {
	es, err = esv7.NewClient(conf)
	if err != nil {
		return nil, err
	}

	res, err := es.Info()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = res.Body.Close()
	}()

	return es, nil
}
