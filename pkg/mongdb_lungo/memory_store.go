package mongdb_lungo

import (
	"context"
	"log"
	"time"

	"get.porter.sh/porter/pkg/portercontext"
	"get.porter.sh/porter/pkg/storage/plugins"
	"github.com/256dpi/lungo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PorterMemoryStore struct {
	*portercontext.Context

	lungo.MemoryStore
	client     lungo.IClient
	engine     *lungo.Engine
	timeout    time.Duration
	database   lungo.IDatabase
	collection lungo.ICollection
}

// NewMemoryStore creates a new storage engine that uses Lungo In-Memory DB.
func NewMemoryStore(c *portercontext.Context, cfg PluginConfig) *PorterMemoryStore {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10 // default to 10 seconds
	}
	return &PorterMemoryStore{
		Context: c,
		timeout: time.Duration(timeout) * time.Second,
	}
}

func (pms *PorterMemoryStore) Connect() (err error) {
	if pms.client != nil {
		return nil
	}

	opts := lungo.Options{
		Store: lungo.NewMemoryStore(),
	}

	cxt, cancel := context.WithTimeout(context.Background(), pms.timeout)
	defer cancel()

	// open database
	pms.client, pms.engine, err = lungo.Open(cxt, opts)
	if err != nil {
		return err
	}

	pms.database = pms.client.Database("porter")
	pms.client.Connect(cxt)

	return nil
}

func (pms *PorterMemoryStore) Close() error {
	if pms.client != nil {
		cxt, cancel := context.WithTimeout(context.Background(), pms.timeout)
		defer cancel()

		err := pms.client.Disconnect(cxt)
		if err != nil {
			return err
		}
		pms.client = nil
	}
	return nil
}

func (pms *PorterMemoryStore) EnsureIndex(opts plugins.EnsureIndexOptions) error {
	// TODO implement me
	panic("implement me")
}

func (pms *PorterMemoryStore) Aggregate(opts plugins.AggregateOptions) (results []bson.Raw, err error) {
	cxt, cancel := context.WithTimeout(context.Background(), pms.timeout)
	defer cancel()

	c := pms.database.Collection(opts.Collection)
	cur, err := c.Aggregate(cxt, opts.Pipeline)
	if err != nil {
		return nil, err
	}

	for cur.Next(cxt) {
		var elem bson.Raw
		err = cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}

	return results, err
}

func (pms *PorterMemoryStore) Count(opts plugins.CountOptions) (int64, error) {
	cxt, cancel := context.WithTimeout(context.Background(), pms.timeout)
	defer cancel()

	c := pms.database.Collection(opts.Collection)
	return c.CountDocuments(cxt, opts.Filter)
}

func (pms *PorterMemoryStore) Find(opts plugins.FindOptions) (results []bson.Raw, err error) {
	cxt, cancel := context.WithTimeout(context.Background(), pms.timeout)
	defer cancel()

	c := pms.database.Collection(opts.Collection)
	findOpts := pms.buildFindOptions(opts)
	cur, err := c.Find(cxt, opts.Filter, findOpts)
	if err != nil {
		return nil, errors.Wrapf(err, "find failed:\n%#v\n%#v", opts.Filter, findOpts)
	}

	// var results []bson.Raw
	// for cur.Next(cxt) {
	// 	results = append(results, cur.Current)
	// }

	for cur.Next(cxt) {
		var elem bson.Raw
		err = cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}

	return results, nil
}

func (pms PorterMemoryStore) buildFindOptions(opts plugins.FindOptions) *options.FindOptions {
	query := options.Find()

	if opts.Select != nil {
		query.SetProjection(opts.Select)
	}

	if opts.Limit > 0 {
		query.SetLimit(opts.Limit)
	}

	if opts.Skip > 0 {
		query.SetSkip(opts.Skip)
	}

	if opts.Sort != nil {
		query.SetSort(opts.Sort)
	}

	return query
}

func (pms PorterMemoryStore) Insert(opts plugins.InsertOptions) error {
	cxt, cancel := context.WithTimeout(context.Background(), pms.timeout)
	defer cancel()

	c := pms.database.Collection(opts.Collection)
	_, err := c.InsertMany(cxt, opts.Documents)
	return err
}

func (pms PorterMemoryStore) Patch(opts plugins.PatchOptions) error {
	cxt, cancel := context.WithTimeout(context.Background(), pms.timeout)
	defer cancel()

	c := pms.database.Collection(opts.Collection)
	_, err := c.UpdateOne(cxt, opts.QueryDocument, opts.Transformation)
	return err
}

func (pms PorterMemoryStore) Remove(opts plugins.RemoveOptions) error {
	cxt, cancel := context.WithTimeout(context.Background(), pms.timeout)
	defer cancel()

	c := pms.database.Collection(opts.Collection)
	if opts.All {
		_, err := c.DeleteMany(cxt, opts.Filter)
		return err
	}
	_, err := c.DeleteOne(cxt, opts.Filter)
	return err
}

func (pms PorterMemoryStore) Update(opts plugins.UpdateOptions) error {
	cxt, cancel := context.WithTimeout(context.Background(), pms.timeout)
	defer cancel()

	c := pms.database.Collection(opts.Collection)
	_, err := c.ReplaceOne(cxt, opts.Filter, opts.Document, &options.ReplaceOptions{Upsert: &opts.Upsert})
	return err
}
