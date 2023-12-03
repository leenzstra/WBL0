package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type K = string
type V = string

type MemoryCacheSuite struct {
	suite.Suite
	cache *MemoryCache[K, V]
	ttl time.Duration
}

func (suite *MemoryCacheSuite) SetupTest() {
	suite.ttl = 5 * time.Minute
	suite.cache = New[K, V](suite.ttl)
}

func (suite *MemoryCacheSuite) TestNewMemoryCacheWithTTL() {
	suite.NotNil(suite.cache)
	suite.IsType(&MemoryCache[K, V]{}, suite.cache)
}

func (suite *MemoryCacheSuite) TestGetItemExists() {
	key := "key"
	value := "value"

	result := suite.cache.SetItem(key, value)
	suite.True(result)

	item, exists := suite.cache.GetItem(key)
	suite.True(exists)
	suite.NotNil(item)

	suite.Equal(value, item)
}

func (suite *MemoryCacheSuite) TestGetItemNotExists() {
	key := "key"

	item, exists := suite.cache.GetItem(key)
	suite.False(exists)
	suite.Empty(item)
}

func (suite *MemoryCacheSuite) TestItemTTL() {
	key := "key"
	value := "value"

	suite.cache.SetItem(key, value)

	result := suite.cache.Cache.Get(key)
	suite.Equal(suite.ttl, result.TTL())
}

func (suite *MemoryCacheSuite) TestHasItem() {
	key := "key"
	value := "value"

	exists := suite.cache.HasItem(key)
	suite.False(exists)

	exists = suite.cache.SetItem(key, value)
	suite.True(exists)
}

func (suite *MemoryCacheSuite) TestLenItems() {
	key := "key"
	value := "value"

	// повторяющийся ключ
	key1 := "key"
	value1 := "value1"

	key2 := "key2"
	value2 := "value2"

	suite.Equal(0, suite.cache.LenItems())
	suite.cache.SetItem(key, value)
	suite.cache.SetItem(key1, value1)
	suite.cache.SetItem(key2, value2)

	suite.Equal(2, suite.cache.LenItems())
}


func TestMemoryCacheSuite(t *testing.T) {
	suite.Run(t, new(MemoryCacheSuite))
}