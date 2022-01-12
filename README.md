# go-ext
[![Go Reference](https://pkg.go.dev/badge/github.com/ShiningRush/goext.svg)](https://pkg.go.dev/github.com/ShiningRush/goext)
[![Go Report Card](https://goreportcard.com/badge/github.com/shiningrush/goext)](https://goreportcard.com/report/github.com/shiningrush/goext)


There are some useful functions and data-structure in development of golang.

## datax
This package contains some extension data structure.

### Set
Descried in [wiki](https://en.wikipedia.org/wiki/Set_(abstract_data_type))
Usage:
```go
    s := datax.NewSet().Add("1", "2", "3", 4)
    s.Len() // 4
    s.Has("3") // true
    s.Remove("3")
    s.TryAdd(4) // false
    s.All() // []interface{}["1", "2", 4]
    
    ano := NewSet().Add("2")
    s.Intersect(ano) // Set["2"]
    s.Union(ano)
    ano.IsSubsetOf(s) // true
```
## parallel
Here are some useful functions help you handle tasks parallel.
### Normal
```go
    parallel.Do(func(workerIdx int) error {
		...
        return nil
    })
```

### Stream parallelism

## stringx
Pleas refer to [gopkg](https://pkg.go.dev/github.com/ShiningRush/goext@v0.0.1/stringx)