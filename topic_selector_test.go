package mercure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	tss, _ := NewTopicSelectorStore()

	assert.True(t, tss.match("https://example.com/foo/bar", "https://example.com/{foo}/bar"))
	_, found := tss.Get("m_https://example.com/{foo}/bar_https://example.com/foo/bar")
	assert.True(t, found)
	assert.True(t, tss.match("https://example.com/foo/bar", "https://example.com/{foo}/bar"))
	assert.False(t, tss.match("https://example.com/foo/bar/baz", "https://example.com/{foo}/bar"))
	//assert.NotNil(t, tss.m["https://example.com/{foo}/bar"].regexp)
	//assert.True(t, tss.m["https://example.com/{foo}/bar"].matchCache["https://example.com/foo/bar"])
	//assert.False(t, tss.m["https://example.com/{foo}/bar"].matchCache["https://example.com/foo/bar/baz"])
	//assert.Equal(t, tss.m["https://example.com/{foo}/bar"].counter, uint32(1))

	assert.True(t, tss.match("https://example.com/kevin/dunglas", "https://example.com/{fistname}/{lastname}"))
	assert.True(t, tss.match("https://example.com/foo/bar", "*"))
	assert.True(t, tss.match("https://example.com/foo/bar", "https://example.com/foo/bar"))
	assert.True(t, tss.match("foo", "foo"))
	assert.False(t, tss.match("foo", "bar"))
}
