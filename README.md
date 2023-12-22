freshcache is an in-memory key:value store/cache with time-based evictions.

It is suitable for applications running on a single machine. It's essentially a
thread-safe map with expiration times. Any object can be stored, for a given
duration or forever, and the cache can be safely used by multiple goroutines.

Although freshcache isn't meant to be used as a persistent datastore, the contents
can be saved to and loaded from a file (using `c.Items()` to retrieve the items
map to serialize, and `NewFrom()` to create a cache from a deserialized one) to
recover from downtime quickly.

The canonical import path is `github.com/vpineda1996/freshcache/v2`, or `github.com/vpineda1996/freshcache` for the v1.
Reference docs are at https://godocs.io/github.com/vpineda1996/freshcache/v2 and
https://godocs.io/github.com/vpineda1996/freshcache

This is a fork of https://github.com/patrickmn/go-cache – which no longer seems
actively maintained. There are two versions of freshcache, both of which are
maintained:



Usage
-----
Some examples from `example_test.go`:

```go
func ExampleSimple() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes.
	//
	// This creates a cache with string keys and values, with Go 1.18 type
	// parameters.
	c := freshcache.New[string, string](5*time.Minute, 10*time.Minute)

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

func ExampleStruct() {
	type MyStruct struct{ Value string }

	// Create a new cache that stores a specific struct.
	c := freshcache.New[string, *MyStruct](freshcache.NoExpiration, freshcache.NoExpiration)
	c.Set("cache", &MyStruct{Value: "value"})

	v, _ := c.Get("cache")
	fmt.Printf("%#v\n", v)

	// Output: &freshcache_test.MyStruct{Value:"value"}
}

func ExampleAny() {
	// Create a new cache that stores any value, behaving similar to freshcache v1
	// or go-cache.
	c := freshcache.New[string, any](freshcache.NoExpiration, freshcache.NoExpiration)

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
	c := freshcache.New[int, *Site](freshcache.NoExpiration, freshcache.NoExpiration)
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
```

Changes
-------

### Compatible changes from go-cache
All these changes are in both:

- Add `Keys()` to list all keys.
- Add `Touch()` to update the expiry on an item.
- Add `GetStale()` to get items even after they've expired.
- Add `Pop()` to get an item and delete it.
- Add `Modify()` to atomically modify existing cache entries (e.g. lists, maps).
- Add `DeleteAll()` to remove all items from the cache with onEvicted call.
- Add `DeleteFunc()` to remove specific items from the cache atomically.
- Add `Rename()` to rename keys, retaining the value and expiry.
- Add `Proxy` type, to access cache items under a different key.
- Various small internal and documentation improvements.

See [issue-list.markdown](/issue-list.markdown) for a complete run-down of the
PRs/issues for go-cache and what was and wasn't included.

FAQ
---

### How can I limit the size of the cache? Is there an option for this?
Not really; freshcache is intended as a thread-safe map with time-based eviction.
This keeps it nice and simple. Adding something like a LRU eviction mechanism
not only makes the code more complex, it also makes the library worse for cases
where you just want a map since it requires additional memory and makes some
operations more expensive (unless a new API is added which make the API worse
for those use cases).

So unless I or someone else comes up with a way to do this which doesn't detract
anything from the simple map use case, I'd rather not add it. Perhaps wrapping
`freshcache.Cache` and overriding some methods could work, but I haven't looked at
it.

tl;dr: this isn't designed to solve every caching use case. That's a feature.
