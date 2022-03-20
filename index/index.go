package index

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"

	"github.com/SecurityBrewery/catalyst/generated/model"
)

type Index struct {
	name     string
	internal bleve.Index
}

func New(name string) (*Index, error) {
	var err error
	var bleveIndex bleve.Index
	if _, oerr := os.Stat(name); os.IsNotExist(oerr) {
		bleveIndex, err = bleve.New(name, bleve.NewIndexMapping())
	} else {
		bleveIndex, err = bleve.Open(name)
	}
	if err != nil {
		return nil, err
	}

	return &Index{name: name, internal: bleveIndex}, nil
}

func (i *Index) Index(incidents []*model.TicketSimpleResponse) {
	b := i.internal.NewBatch()
	for _, incident := range incidents {
		if incident.ID == 0 {
			log.Println(errors.New("no ID"), incident)

			continue
		}

		err := b.Index(fmt.Sprint(incident.ID), incident)
		if err != nil {
			log.Println(err)
		}
	}

	if err := i.internal.Batch(b); err != nil {
		log.Println(err)
	}
}

func (i *Index) Search(term string) (ids []string, err error) {
	query := bleve.NewQueryStringQuery(term)
	result, err := i.internal.Search(bleve.NewSearchRequestOptions(query, 10000, 0, false))
	if err != nil {
		return nil, err
	}
	for _, match := range result.Hits {
		ids = append(ids, match.ID)
	}

	return ids, nil
}

func (i *Index) Truncate() error {
	err := i.internal.Close()
	if err != nil {
		return err
	}
	err = os.RemoveAll(i.name)
	if err != nil {
		return err
	}
	index, err := bleve.New(i.name, bleve.NewIndexMapping())
	if err != nil {
		return err
	}
	i.internal = index

	return nil
}

func (i *Index) Close() error {
	return i.internal.Close()
}
