package mercure

import (
	"regexp"
	"strings"

	"github.com/dgraph-io/ristretto"
	uritemplate "github.com/yosida95/uritemplate/v3"
)

// topicSelectorStore caches compiled templates to improve memory and CPU usage.
type TopicSelectorStore struct {
	*ristretto.Cache
}

// NewTopicSelectorStore creates a new topic selector store.
func NewTopicSelectorStore() (*TopicSelectorStore, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err == nil {
		return &TopicSelectorStore{cache}, nil
	}

	return nil, err
}

func (tss *TopicSelectorStore) match(topic, topicSelector string, addToCache bool) bool {
	// Always do an exact matching comparison first
	// Also check if the topic selector is the reserved keyword *
	if topicSelector == "*" || topic == topicSelector {
		return true
	}

	r := tss.getRegexp(topicSelector, addToCache)
	if r == nil {
		return false
	}

	k := "m" + "_" + topicSelector + "_" + topic
	value, found := tss.Get(k)
	if found {
		return value.(bool)
	}

	// Use template.Regexp() instead of template.Match() for performance
	// See https://github.com/yosida95/uritemplate/pull/7
	match := r.MatchString(topic)
	tss.Set(k, match, 1)

	return match
}

// getRegexp retrieves regexp for this template selector.
func (tss *TopicSelectorStore) getRegexp(topicSelector string, addToCache bool) *regexp.Regexp {
	// If it's definitely not an URI template, skip to save some resources
	if !strings.Contains(topicSelector, "{") {
		return nil
	}

	k := "t" + topicSelector

	value, found := tss.Get(k)
	if found {
		return value.(*regexp.Regexp)
	}

	// If an error occurs, it's a raw string
	if tpl, err := uritemplate.New(topicSelector); err == nil {
		r := tpl.Regexp()
		tss.Set(k, r, 10)

		return r
	}

	return nil
}

// cleanup removes unused compiled templates from memory.
func (tss *TopicSelectorStore) cleanup(topics []string) {
	// FIXME: to remove
}
