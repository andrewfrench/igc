package cache

import (
	"afren.ch/set"
	"afren.ch/env"
	"time"
	"log"
)

type setCache struct {
	sets     map[string]*set.CorrelationSet
	keyQueue []queueElement
	maxSize  int
	maxAge   time.Duration
}

type queueElement struct {
	key string
	added time.Time
}

var c *setCache

func init() {
	c = new(setCache)
	c.sets = map[string]*set.CorrelationSet{}
	c.maxSize = env.OptionalInt("CACHE_SIZE", 256)
	c.maxAge = time.Duration(env.OptionalInt("CACHE_MAX_AGE_MIN", 10)) * time.Minute

	go func() {
		for range time.Tick(10 * time.Minute) {
			removeExpired()
			reportSize()
		}
	}()
}

func Contains(base string) bool {
	_, exists := c.sets[base]

	return exists
}

func Insert(set *set.CorrelationSet) {
	if len(c.keyQueue) == c.maxSize {
		delete(c.sets, c.keyQueue[0].key)	// The first element is the oldest
		c.keyQueue = c.keyQueue[1:]
	}

	c.keyQueue = append(c.keyQueue, queueElement{
		key: set.GetBase(),
		added: time.Now(),
	})
	c.sets[set.GetBase()] = set
}

func Count() int {
	return len(c.sets)
}

func Get(base string) *set.CorrelationSet {
	return c.sets[base]
}

func removeExpired() {
	if len(c.keyQueue) == 0 { return }
	if time.Since(c.keyQueue[0].added) > c.maxAge {
		key := c.keyQueue[0].key
		log.Printf("%s expired, removing from cache", key)
		delete(c.sets, key)

		if len(c.keyQueue) > 1 {
			c.keyQueue = c.keyQueue[1:]
		} else {
			c.keyQueue = []queueElement{}
		}

		removeExpired()
	}
}

func reportSize() {
	size := len(c.sets)

	if size > 0 { log.Printf("Cache size: %d", len(c.sets)) }
}