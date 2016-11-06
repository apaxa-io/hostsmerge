package main

import (
	"encoding/json"
	"os"
	"path"
	"sync"
	"time"
)

type cacheElement struct {
	Exists bool
	Till   time.Time
}

type dnsCache struct {
	m map[string]cacheElement
	l sync.RWMutex
}

func (c *dnsCache) CleanOutdated() (removed int) {
	now := time.Now()
	c.l.Lock()
	defer c.l.Unlock()
	for key, value := range c.m {
		if value.Till.Before(now) {
			delete(c.m, key)
			removed++
		}
	}
	return
}

func (c *dnsCache) Add(fqdn string, exists bool, till time.Time) {
	if till.Before(time.Now()) {
		return
	}

	c.l.Lock()
	defer c.l.Unlock()
	c.m[fqdn] = cacheElement{Exists: exists, Till: till}
}

func (c *dnsCache) Get(fqdn string) (inCache, exists bool) {
	c.l.RLock()
	defer c.l.RUnlock()
	t, inCache := c.m[fqdn]
	if inCache {
		exists = t.Exists
	}
	return
}

func (c *dnsCache) Len() int {
	c.l.RLock()
	defer c.l.RUnlock()
	return len(c.m)
}

func (c *dnsCache) ReadFromFile(fn string) error {
	c.l.Lock()
	defer c.l.Unlock()

	c.m = make(map[string]cacheElement)

	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&c.m)
}

func (c *dnsCache) WriteToFile(fn string) error {
	c.l.Lock()
	defer c.l.Unlock()

	if err := os.MkdirAll(path.Dir(fn), 0777); err != nil {
		return err
	}

	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(c.m)
}
