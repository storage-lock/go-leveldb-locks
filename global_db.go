package leveldb_locks

import (
	"context"
	leveldb_storage "github.com/storage-lock/go-leveldb-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"github.com/syndtr/goleveldb/leveldb"
)

var dbStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[*leveldb.DB, *leveldb.DB] = storage_lock_factory.NewStorageLockFactoryBeanFactory[*leveldb.DB, *leveldb.DB]()

func NewLevelDBLockByDB(ctx context.Context, db *leveldb.DB, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetLevelDBLockFactoryByDB(ctx, db)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewLevelDBLockByDBWithOptions(ctx context.Context, db *leveldb.DB, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetLevelDBLockFactoryByDB(ctx, db)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

func GetLevelDBLockFactoryByDB(ctx context.Context, db *leveldb.DB) (*storage_lock_factory.StorageLockFactory[*leveldb.DB], error) {
	return dbStorageLockFactoryBeanFactory.GetOrInit(ctx, db, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*leveldb.DB], error) {
		connectionManager := leveldb_storage.NewLevelDBConnectionManagerFromDB(db)
		options := leveldb_storage.NewLevelDBStorageOptions().SetConnectionManager(connectionManager)
		storage, err := newLevelDBStorage(ctx, options)
		if err != nil {
			return nil, err
		}
		factory := storage_lock_factory.NewStorageLockFactory(storage, options.ConnectionManager)
		return factory, nil
	})
}
