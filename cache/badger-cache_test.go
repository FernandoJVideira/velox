package cache

import "testing"

func TestBadgerCache_Has(t *testing.T) {
	err := testBadgerCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo should not exist in the cache")
	}

	err = testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	has, err := testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !has {
		t.Error("foo should exist in the cache")
	}

	_ = testBadgerCache.Forget("foo")
}

func TestBadgerCache_Get(t *testing.T) {
	err := testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	x, err := testBadgerCache.Get("foo")
	if err != nil {
		t.Error(err)
	}

	if x != "bar" {
		t.Error("did not get the expected value")
	}
}

func TestBadgerCache_Forget(t *testing.T) {
	err := testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	has, err := testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if has {
		t.Error("foo should not exist in the cache")
	}
}

func TestBadgerCache_EmptyByMatch(t *testing.T) {
	err := testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Set("foo:bar", "baz")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.EmptyByMatch("foo")
	if err != nil {
		t.Error(err)
	}

	has, err := testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if has {
		t.Error("foo should not exist in the cache")
	}

	has, err = testBadgerCache.Has("foo:bar")
	if err != nil {
		t.Error(err)
	}

	if has {
		t.Error("foo:bar should not exist in the cache")
	}

	has, err = testBadgerCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if !has {
		t.Error("alpha should exist in the cache")
	}
}

func TestBadgerCache_Empty(t *testing.T) {
	err := testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Set("foo:bar", "baz")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Empty()
	if err != nil {
		t.Error(err)
	}

	has, err := testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if has {
		t.Error("foo should not exist in the cache")
	}

	has, err = testBadgerCache.Has("foo:bar")
	if err != nil {
		t.Error(err)
	}

	if has {
		t.Error("foo:bar should not exist in the cache")
	}

	has, err = testBadgerCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}

	if has {
		t.Error("alpha should not exist in the cache")
	}
}
