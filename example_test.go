package freshcache_test

import (
	"fmt"
	"time"

	"github.com/vpineda1996/freshcache"
)

func ExampleCache_simple() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes.
	//
	// This creates a cache with string keys and values, with Go 1.18 type
	// parameters.
	c := freshcache.New[string, string](100, 5*time.Minute, 10*time.Minute)

	// Set the value of the key "foo" to "bar", with the default expiration.
	c.Set("foo", "bar")

	// Set the value of the key "baz" to "never", with no expiration time. The
	// item won't be removed until it's removed with c.Delete("baz").
	c.SetWithExpire("baz", "never", freshcache.NoExpiration)

	// Get the value associated with the key "foo" from the cache; due to the
	// use of type parameters this is a string, and no type assertions are
	// needed.
	foo, ok := c.Get("foo")
	if ok {
		fmt.Println(foo)
	}

	// Output: bar
}

func ExampleCache_struct() {
	type MyStruct struct{ Value string }

	// Create a new cache that stores a specific struct.
	c := freshcache.New[string, *MyStruct](100, freshcache.NoExpiration, freshcache.NoExpiration)
	c.Set("cache", &MyStruct{Value: "value"})

	v, _ := c.Get("cache")
	fmt.Printf("%#v\n", v)

	// Output: &freshcache_test.MyStruct{Value:"value"}
}

func ExampleCache_any() {
	// Create a new cache that stores any value, behaving similar to freshcache v1
	// or go-cache.
	c := freshcache.New[string, any](100, freshcache.NoExpiration, freshcache.NoExpiration)

	c.Set("a", "value 1")
	c.Set("b", 42)

	a, _ := c.Get("a")
	b, _ := c.Get("b")

	// This needs type assertions.
	p := func(a string, b int) { fmt.Println(a, b) }
	p(a.(string), b.(int))

	// Output: value 1 42
}

func ExampleProxy() {
	type Site struct {
		ID       int
		Hostname string
	}

	site := &Site{
		ID:       42,
		Hostname: "example.com",
	}

	// Create a new site which caches by site ID (int), and a "proxy" which
	// caches by the hostname (string).
	c := freshcache.New[int, *Site](100, freshcache.NoExpiration, freshcache.NoExpiration)
	p := freshcache.NewProxy[string, int, *Site](c)

	p.Set(42, "example.com", site)

	siteByID, ok := c.Get(42)
	fmt.Printf("%v %v\n", ok, siteByID)

	siteByHost, ok := p.Get("example.com")
	fmt.Printf("%v %v\n", ok, siteByHost)

	// They're both the same object/pointer.
	fmt.Printf("%v\n", siteByID == siteByHost)

	// Output:
	// true &{42 example.com}
	// true &{42 example.com}
	// true
}
