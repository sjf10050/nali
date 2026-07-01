package db

import (
	"sync"

	"github.com/zu1k/nali/pkg/dbif"
)

var (
	// cacheMu guards dbNameCache and dbTypeCache (plain maps). queryCache is a
	// sync.Map and needs no external locking.
	cacheMu     sync.RWMutex
	dbNameCache = make(map[string]dbif.DB)
	dbTypeCache = make(map[dbif.QueryType]dbif.DB)
	queryCache  = sync.Map{}
)

var (
	NameDBMap = make(NameMap)
	TypeDBMap = make(TypeMap)
)
