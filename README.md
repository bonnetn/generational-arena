# Generational area

Just a simple implementation of a *generational arena*.

A generational arena is a container. Callers can insert, remove and get values, similar to arrays.
However, the order is not guaranteed.

All operations run in constant-time complexity.

Under the hood, the values are stored in an array, which provides good locality of reference.
The library attempts to re-use the empty slots in the array, so the memory footprint is kept low.

### Example:
```go
import "github.com/bonnetn/generational-arena/arena"

var a arena.Arena[string]
idHello := a.Insert("hello")
idWorld := a.Insert("world")

a.Remove(idHello)
a.Insert("hi, ")

// This will print "world"
fmt.Println(a.Get(idWorld)

// Iterate over the contents:
it := a.Iterate()
for {
	v, _, done := a.Next()
	if done {
        break
    }
    // This will print "hi" and "world".
    fmt.Println(v)
}
```

### How it works

The values are stored in an array.
Similar to a slab, the arena attempts to insert the new values in the array where there is a "hole".

The identifier contains both the index of the value in the array and the generation of the value.
The *generation* is a counter that is incremented every time a value is removed.
This way, the arena can detect if the value has been removed and the index is now occupied by a new value.


#### Example:
```
Initial array contains two elements:
[ {value="hello", generation=1}, {value="world", generation=1} ]

The first element gets remove, there is now a hole:
[ {EMPTY, generation=1}, {value="world", generation=1} ]

Upon insertion, the first empty cell will be populated (notice that the generation count is incremented):
[ {value="hi", generation=2}, {value="world", generation=1} ]

If no space is left, the new element is appended at the end:
[ {value="hi", generation=2}, {value="world", generation=1}, {value="wow", generation=1} ]
```

### Zero allocation variant

The library provides a zero-allocation variant, which is useful for high-performance applications.

```go
import "github.com/bonnetn/generational-arena/arena"

var ( // Pre-allocate the underlying arrays
    buffer   [1024]int
    metadata [1024]MetadataEntry
)
a := arena.NewFixedSizeArena[int](buffer[:], metadata[:])

// The rest of the code is the same.
// The arena will not allocate any memory.
// The insertions will panic if the buffer is full.
```


 


