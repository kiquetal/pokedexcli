package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key   string
		value []byte
	}{
		{"key1", []byte("some text response")},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.value)

			value, found := cache.Get(c.key)
			fmt.Printf("value: %v\n", string(value))
			if !found {
				t.Errorf("expected key %q to be found", c.key)
			}
			if string(value) != string(c.value) {
				t.Errorf("expected value %q, got %q", c.value, value)
			}
		})

	}
}
