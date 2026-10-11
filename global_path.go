package leveldb_locks

import (
	"context"
	leveldb_storage "github.com/storage-lock/go-leveldb-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"github.com/syndtr/goleveldb/leveldb"
)

var pathStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[string, *leveldb.DB] = storage_lock_factory.NewStorageLockFactoryBeanFactory[string, *leveldb.DB]()

var newLevelDBStorage = leveldb_storage.NewLevelDBStorage

func NewLevelDBLockByPath(ctx context.Context, path string, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetLevelDBLockFactoryByPath(ctx, path)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewLevelDBLockByPathWithOptions(ctx context.Context, path string, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetLevelDBLockFactoryByPath(ctx, path)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

func GetLevelDBLockFactoryByPath(ctx context.Context, path string) (*storage_lock_factory.StorageLockFactory[*leveldb.DB], error) {
	return pathStorageLockFactoryBeanFactory.GetOrInit(ctx, path, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*leveldb.DB], error) {
		connectionManager := leveldb_storage.NewLevelDBConnectionManagerFromPath(path)
		options := leveldb_storage.NewLevelDBStorageOptions().SetConnectionManager(connectionManager)
		storage, err := newLevelDBStorage(ctx, options)
		if err != nil {
			return nil, err
		}
		factory := storage_lock_factory.NewStorageLockFactory(storage, options.ConnectionManager)
		return factory, nil
	})
}
