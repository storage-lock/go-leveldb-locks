package leveldb_locks

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewLevelDBLockByPath(t *testing.T) {
	// 创建临时目录作为 LevelDB 数据库路径
	tmpDir, err := os.MkdirTemp("", "storage-lock-leveldb-locks-test-*")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpDir)

	ctx := context.Background()
	lockId := "test-leveldb-lock"

	factory, err := GetLevelDBLockFactoryByPath(ctx, tmpDir)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	lock, err := NewLevelDBLockByPath(ctx, tmpDir, lockId)
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	// 加锁
	ownerId := "test-owner"
	err = lock.Lock(ctx, ownerId)
	assert.Nil(t, err)

	// 解锁
	err = lock.UnLock(ctx, ownerId)
	assert.Nil(t, err)
}
