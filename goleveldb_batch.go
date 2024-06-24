package db

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type goLevelDBBatch struct {
	db    *GoLevelDB
	batch *leveldb.Batch
}

var _ Batch = (*goLevelDBBatch)(nil)

func newGoLevelDBBatch(db *GoLevelDB) *goLevelDBBatch {
	return &goLevelDBBatch{
		db:    db,
		batch: new(leveldb.Batch),
	}
}

// Set implements Batch.
func (b *goLevelDBBatch) Set(key, value []byte) error {
	// log.Printf("Set (batch): name is %s, key is %s, value is %d bytes", b.db.name, hex.EncodeToString(key), len(value))
	if len(key) == 0 {
		return errKeyEmpty
	}
	if value == nil {
		return errValueNil
	}
	if b.batch == nil {
		return errBatchClosed
	}
	b.batch.Put(key, value)
	return nil
}

// Delete implements Batch.
func (b *goLevelDBBatch) Delete(key []byte) error {
	if len(key) == 0 {
		return errKeyEmpty
	}
	if b.batch == nil {
		return errBatchClosed
	}
	b.batch.Delete(key)
	return nil
}

// Write implements Batch.
func (b *goLevelDBBatch) Write() error {
	return b.write(false)
}

// WriteSync implements Batch.
func (b *goLevelDBBatch) WriteSync() error {
	return b.write(true)
}

func (b *goLevelDBBatch) write(sync bool) error {
	log.Printf("Write (batch): name is %s, size is %d bytes", b.db.name, len(b.batch.Dump()))
	if b.batch == nil {
		return errBatchClosed
	}
	err := b.db.db.Write(b.batch, &opt.WriteOptions{Sync: sync})
	if err != nil {
		return err
	}
	// Make sure batch cannot be used afterwards. Callers should still call Close(), for errors.
	return b.Close()
}

// Close implements Batch.
func (b *goLevelDBBatch) Close() error {
	if b.batch != nil {
		b.batch.Reset()
		b.batch = nil
	}
	return nil
}
