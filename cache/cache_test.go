package cache

import "testing"

func TestRedisCacheHas(t *testing.T) {
	err := testRedisCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo should not be in the cache")
	}

	err = testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	inCache, err = testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("foo should be in the cache")
	}
}

func TestRedisCacheGet(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	value, err := testRedisCache.Get("foo")
	if err != nil {
		t.Error(err)
	}

	if value != "bar" {
		t.Error("did not get the expected value")
	}
}

func TestRedisCacheForget(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo should not be in the cache")
	}
}

func TestRedisCacheEmpty(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Empty()
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo should not be in the cache")
	}

	inCache, err = testRedisCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foobar should not be in the cache")
	}
}

func TestRedisCacheEmptyByMatch(t *testing.T) {
	err := testRedisCache.Set("alpha", "foo")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("alpha2", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("beta", "baz")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.EmptyByMatch("alpha")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("alpha should not be in the cache")
	}

	inCache, err = testRedisCache.Has("alpha2")
	if err != nil {
		t.Error("alpha2 should not be in the cache")
	}

	if inCache {
		t.Error("alpha2 should not be in the cache")
	}

	inCache, err = testRedisCache.Has("beta")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("beta should be in the cache")
	}
}

func TestEncodeDecode(t *testing.T) {
	item := Entry{}
	item["foo"] = "bar"

	bytes, err := encode(item)
	if err != nil {
		t.Error(err)
	}

	_, err = decode(string(bytes))
	if err != nil {
		t.Error(err)
	}
}
