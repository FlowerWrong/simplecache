package simplecache

import (
	"testing"
	"time"
)

var c *CacheMap

var (
	key1 = "key1"
	val1 = "val1"
)

func init() {
	c = New(1)
	c.SetMaxMemory("1MB")
}

func TestNew(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	if c.nbytes != 0 {
		t.Error(c.nbytes)
	}
}

func TestSetMaxMemory(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	c2 := New(1)
	c2.SetMaxMemory("1B")
	err := c2.Set(key1, val1, time.Second)
	if err == nil {
		t.Error("Max memory limit not work")
	}
}

func TestSet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	err := c.Set(key1, val1, time.Second)
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	c.Flush()
	_, ok := c.Get(key1)
	if ok {
		t.Error("God create it")
	}

	err := c.Set(key1, val1, time.Second)
	if err != nil {
		t.Error(err)
	}

	item, ok := c.Get(key1)
	if !ok {
		t.Error("Not found")
	}
	if item.(string) != val1 {
		t.Error("not equal")
	}
}

func TestDel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	c.Flush()
	err := c.Set(key1, val1, time.Second)
	if err != nil {
		t.Error(err)
	}

	item, ok := c.Get(key1)
	if !ok {
		t.Error("Not found")
	}
	if item.(string) != val1 {
		t.Error("not equal")
	}

	c.Del(key1)
	_, ok = c.Get(key1)
	if ok {
		t.Error("God create it")
	}
}

func TestExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	c.Flush()
	err := c.Set(key1, val1, time.Second)
	if err != nil {
		t.Error(err)
	}

	if !c.Exists(key1) {
		t.Error("not equal")
	}
}

func TestSize(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	c.Flush()
	err := c.Set(key1, val1, time.Second)
	if err != nil {
		t.Error(err)
	}

	if c.Size() != 1 {
		t.Error("only 1 key")
	}
}

func TestKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	c.Flush()
	err := c.Set(key1, val1, time.Second)
	if err != nil {
		t.Error(err)
	}

	if len(c.Keys()) != 1 {
		t.Error("not equal")
	}
	if c.Keys()[0] != key1 {
		t.Error("not equal")
	}
}

func TestGC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	c.Flush()
	err := c.Set(key1, val1, time.Second)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 1)

	if c.Exists(key1) {
		t.Error("GC not work")
	}
}
