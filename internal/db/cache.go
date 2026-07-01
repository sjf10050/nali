package db

import (
	"container/list"
	"sync"

	"github.com/zu1k/nali/pkg/dbif"
)

var (
	cacheMu     sync.RWMutex
	dbNameCache = make(map[string]dbif.DB)
	dbTypeCache = make(map[dbif.QueryType]dbif.DB)
	queryCache  = newBoundedCache(10000)
)

var (
	NameDBMap = make(NameMap)
	TypeDBMap = make(TypeMap)
)

// boundedCache is a simple LRU cache for query results backed by
// container/list. It is safe for concurrent use.
type boundedCache struct {
	mu      sync.Mutex
	items   map[string]*list.Element
	lruList *list.List
	maxSize int
}

type cacheEntry struct {
	key   string
	value *Result
}

func newBoundedCache(maxSize int) *boundedCache {
	return &boundedCache{
		items:   make(map[string]*list.Element),
		lruList: list.New(),
		maxSize: maxSize,
	}
}

// Load returns the cached result and true if the key is present.
// The entry is promoted to the front of the LRU list.
func (c *boundedCache) Load(key string) (*Result, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[key]; ok {
		c.lruList.MoveToFront(elem)
		return elem.Value.(*cacheEntry).value, true
	}
	return nil, false
}

// Store inserts or updates a key, evicting the least recently used entry
// when the cache exceeds maxSize.
func (c *boundedCache) Store(key string, value *Result) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[key]; ok {
		c.lruList.MoveToFront(elem)
		elem.Value.(*cacheEntry).value = value
		return
	}
	for c.lruList.Len() >= c.maxSize {
		oldest := c.lruList.Back()
		if oldest != nil {
			c.lruList.Remove(oldest)
			delete(c.items, oldest.Value.(*cacheEntry).key)
		}
	}
	entry := &cacheEntry{key: key, value: value}
	elem := c.lruList.PushFront(entry)
	c.items[key] = elem
}
