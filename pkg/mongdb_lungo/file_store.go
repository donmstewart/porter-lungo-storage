package mongdb_lungo

import (
	"time"

	"get.porter.sh/porter/pkg/portercontext"
	"get.porter.sh/porter/pkg/storage/plugins"
	"github.com/256dpi/lungo"
	"go.mongodb.org/mongo-driver/bson"
)

type PorterFileStore struct {
	*portercontext.Context

	lungo.FileStore
	timeout time.Duration
	path    string
}

// NewFileStore creates a new storage engine that uses Lungo In-Memory DB.
func NewFileStore(c *portercontext.Context, cfg PluginConfig) *PorterFileStore {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10 // default to 10 seconds
	}
	return &PorterFileStore{
		Context: c,
		path:    cfg.Path,
		timeout: time.Duration(timeout) * time.Second,
	}
}

func (p PorterFileStore) Connect() error {
	// TODO implement me
	panic("implement me")
}

func (p PorterFileStore) Close() error {
	// TODO implement me
	panic("implement me")
}

func (p PorterFileStore) EnsureIndex(opts plugins.EnsureIndexOptions) error {
	// TODO implement me
	panic("implement me")
}

func (p PorterFileStore) Aggregate(opts plugins.AggregateOptions) ([]bson.Raw, error) {
	// TODO implement me
	panic("implement me")
}

func (p PorterFileStore) Count(opts plugins.CountOptions) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (p PorterFileStore) Find(opts plugins.FindOptions) ([]bson.Raw, error) {
	// TODO implement me
	panic("implement me")
}

func (p PorterFileStore) Insert(opts plugins.InsertOptions) error {
	// TODO implement me
	panic("implement me")
}

func (p PorterFileStore) Patch(opts plugins.PatchOptions) error {
	// TODO implement me
	panic("implement me")
}

func (p PorterFileStore) Remove(opts plugins.RemoveOptions) error {
	// TODO implement me
	panic("implement me")
}

func (p PorterFileStore) Update(opts plugins.UpdateOptions) error {
	// TODO implement me
	panic("implement me")
}
