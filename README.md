# Simple cache

## Usage

```golang
c := simplecache.New(1)
if !c.SetMaxMemory("2MB") {
    log.Panicln("Set max memory for cache failed")
}

for i := 0; i < 1000; i++ {
    err := c.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i), time.Second)
    if err != nil {
        log.Panicln(err)
    }
}

v, ok := c.Get(fmt.Sprintf("key%d", 1))
if !ok {
    log.Panicln("Not found")
}
log.Println("Found", v.(string))

c.Del(fmt.Sprintf("key%d", 1))

log.Println("keys ", c.Keys())

log.Println("keys len", c.Size())

c.Flush()

time.Sleep(time.Second * 2) // for GC
```
