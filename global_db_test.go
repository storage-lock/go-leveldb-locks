package leveldb_locks

import (
	"context"
	"errors"
	"os"
	"testing"

	leveldb_storage "github.com/storage-lock/go-leveldb-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestNewLevelDBLockByDB(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "storage-lock-leveldb-locks-db-*")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpDir)

	db, err := leveldb.OpenFile(tmpDir, nil)
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.Background()
	factory, err := GetLevelDBLockFactoryByDB(ctx, db)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 缓存命中
	factory2, err := GetLevelDBLockFactoryByDB(ctx, db)
	assert.Nil(t, err)
	assert.Equal(t, factory, factory2)

	lock, err := NewLevelDBLockByDB(ctx, db, "test-lock-db")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	assert.Nil(t, lock.Lock(ctx, "owner-1"))
	assert.Nil(t, lock.UnLock(ctx, "owner-1"))
}

func TestNewLevelDBLockByDBWithOptions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "storage-lock-leveldb-locks-dbopt-*")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpDir)

	db, err := leveldb.OpenFile(tmpDir, nil)
	assert.Nil(t, err)
	defer db.Close()

	options := storage_lock.NewStorageLockOptionsWithLockId("test-lock-db-opt")
	lock, err := NewLevelDBLockByDBWithOptions(context.Background(), db, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock)
}

func TestNewLevelDBLockByPathWithOptions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "storage-lock-leveldb-locks-pathopt-*")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpDir)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-lock-path-opt")
	lock, err := NewLevelDBLockByPathWithOptions(context.Background(), tmpDir, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock)
}

func TestLevelDBLockFactoryByPath_Errors(t *testing.T) {
	old := newLevelDBStorage
	newLevelDBStorage = func(ctx context.Context, options *leveldb_storage.LevelDBStorageOptions) (*leveldb_storage.LevelDBStorage, error) {
		return nil, errors.New("boom")
	}
	defer func() { newLevelDBStorage = old }()

	// ByPath 错误分支
	_, err := GetLevelDBLockFactoryByPath(context.Background(), "/tmp/leveldb-locks-err-path")
	assert.NotNil(t, err)
	_, err = NewLevelDBLockByPath(context.Background(), "/tmp/leveldb-locks-err-path", "x")
	assert.NotNil(t, err)
	_, err = NewLevelDBLockByPathWithOptions(context.Background(), "/tmp/leveldb-locks-err-path", storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)

	// ByDB 错误分支（用独立 db 避免缓存污染）
	tmpDir, err := os.MkdirTemp("", "storage-lock-leveldb-locks-dberr-*")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpDir)
	db, err := leveldb.OpenFile(tmpDir, nil)
	assert.Nil(t, err)
	defer db.Close()
	_, err = GetLevelDBLockFactoryByDB(context.Background(), db)
	assert.NotNil(t, err)
	_, err = NewLevelDBLockByDB(context.Background(), db, "x")
	assert.NotNil(t, err)
	_, err = NewLevelDBLockByDBWithOptions(context.Background(), db, storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)
}
