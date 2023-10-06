package core

import "sync"

type ServiceStats struct {
	ID      string
	Type    string
	Success bool
	EndTime int64
}

type store struct {
	ServiceStats   chan *ServiceStats
	ServiceStats2  []*ServiceStats
	RunningService *sync.Map
}

// TODO: memory leak, need to limit the size of ServiceStats2
var Store = store{
	ServiceStats:   make(chan *ServiceStats),
	ServiceStats2:  make([]*ServiceStats, 0),
	RunningService: &sync.Map{},
}
