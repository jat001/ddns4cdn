package core

import (
	"sync"
	"time"
)

type ServiceStats struct {
	ID      string
	Type    string
	Success bool
	EndTime time.Time
}

type store struct {
	ServiceChan    chan *ServiceStats
	ServiceMap     map[string][]*ServiceStats
	RunningService *sync.Map
}

var Store = store{
	ServiceChan:    make(chan *ServiceStats),
	ServiceMap:     make(map[string][]*ServiceStats),
	RunningService: &sync.Map{},
}
