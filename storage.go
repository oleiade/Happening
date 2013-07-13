package happening

type KvPair struct {
	Key   []byte
	Value []byte
}

type StorageBackend interface {
	Get(key []byte)
	Put(pair KvPair)
	Delete(key []byte)
	MGet(keys [][]byte)
	MPut([]KvPair)
}
