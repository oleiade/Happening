package happening

import (
	"bytes"
	"path/filepath"
	leveldb "github.com/jmhodges/levigo"
)

type LeveldbBackend struct {
	Options *leveldb.Options
	Db      *leveldb.DB
}

// NewLeveldbBackend creates a new leveldb database connector.
func NewLeveldbBackend(storagePath string) (backend *LeveldbBackend, err error) {
	// Set up backend to use a lru cache and
	// create store files if not existing yet
	opts := leveldb.NewOptions()
	opts.SetCache(leveldb.NewLRUCache(LEVELDB_LRU_CACHE_SIZE))
	opts.SetCreateIfMissing(true)

	// Open database file
	db, err := leveldb.Open(filepath.Join(storagePath, "data"), opts)
	if err != nil {
		return nil, err
	}

	// Build backend using previously created
	// options and db connector
	backend = &LeveldbBackend{
		Options: opts,
		Db:      db,
	}

	return backend, nil
}

func (backend *LeveldbBackend) Get(key []byte) (value []byte, err error) {
	ro := leveldb.NewReadOptions()
	value, err = backend.Db.Get(ro, key)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (backend *LeveldbBackend) Put(pair KvPair) error {
	wo := leveldb.NewWriteOptions()
	return backend.Db.Put(wo, pair.Key, pair.Value)
}

func (backend *LeveldbBackend) Delete(key []byte) (err error) {
	wo := leveldb.NewWriteOptions()
	return backend.Db.Delete(wo, key)
}

func (backend *LeveldbBackend) MGet(keys [][]byte) (values [][]byte, err error) {
	var data [][]byte = make([][]byte, len(keys))

	// Read over a Db read-only snapshot
	readOptions := leveldb.NewReadOptions()
	snapshot := backend.Db.NewSnapshot()
	readOptions.SetSnapshot(snapshot)

	if len(keys) > 0 {
		// Extract start -> end key range to extract
		// from leveldb using an iterator
		start := keys[0]
		end := keys[len(keys)-1]

		// Keep track of the input keys index
		// in order to output results in order
		keysIndex := make(map[string]int)
		for index, element := range keys {
			keysIndex[string(element)] = index
		}

		// Fetch values using an iterator
		it := backend.Db.NewIterator(readOptions)
		defer it.Close()
		it.Seek([]byte(start))

		for ; it.Valid(); it.Next() {
			if bytes.Compare(it.Key(), []byte(end)) > 1 {
				break
			}

			if index, present := keysIndex[string(it.Key())]; present {
				data[index] = it.Value()
			}
		}
	}

	backend.Db.ReleaseSnapshot(snapshot)

	return data, nil
}

func (backend *LeveldbBackend) MPut(pairs []KvPair) (err error) {
	var batch *leveldb.WriteBatch = leveldb.NewWriteBatch()

	for _, pair := range pairs {
		batch.Put(pair.Key, pair.Value)
	}

	wo := leveldb.NewWriteOptions()
	return backend.Db.Write(wo, batch)
}
