# seq

Clojure-like immutable data-structures and lazy lists for go.

## Installation

```
    go get github.com/mediocregopher/seq
```

or if you're using [goat][goat]:

```
    - loc: https://github.com/mediocregopher/seq.git
      type: git
      ref: v0.1.0
      path: github.com/mediocregopher/seq
```

## Usage

Check the [godocs][godocs] for the full API. [Examples][examples] are a good
place to look too.

## About

This library constitutes an attempt at bringing immutability and laziness to go
in a thread-safe way, at the cost of type-safety and code-cleanliness.

There are four available types:

* `List` - Single linked list
* `Set` - Hash-tree based unordered set
* `HashMap` - A simple key/value hash map built on top of `Set`
* `Lazy` - Lazily evaluated sequence

All four of these implement the `Seq` interface, which simply provides a way to
iterate over the structure a single time.

### Immutability

All operations on seq's datastructures are immutable, meaning they won't effect
the original variable but instead return a new one with the change in place.
Conceptually a copy is done. In reality seq uses structure sharing so only a
minimal amount of copying is actually necessary.

Since all seq variables are immutable, they are also inherently thread-safe (or
go-routine safe). They can be passed around and operated on by any number of
go-routines with woeful abandon.

### Laziness

The `Lazy` type is a special seq type which is conceptually similar to a
linked-list. Where it differs is that it only evalutates its elements as needed,
so if you only consume half of the "list" only half of the elements will
actually be created.

Seq provides a number of methods which operate on types implementing the `Seq`
interface and return a `Lazy`. These include `LMap`, `LFilter`, and `LTake`.
Since `Lazy` also implements `Seq`, you can chain together lazy functions like
so:

```go
result := seq.LFilter(filterFn, seq.LMap(mapFn, l))
```

In the above example, no intermediate lists are created like if you had used
`Map` or `Filter`. Additionally, `Lazy` will cache its results, so iterating
over `result` multiple times would not cause the `filterFn` and `mapFn`
functions to be called more than once on each element each. This caching is
thread-safe, so multiple go-routines can iterate over the same `Lazy` safely in
all cases.

## Disclaimer

This library has its upsides and downsides, and is probably only truly useful
for a minority of cases. But I had fun making it, and learned a lot, so I
figured maybe someone else would to.

# Legal

The [license][license] supplied in the repo applies to all code and resources
provided in the repo.

Copyright Brian Picciano, 2014

[goat]: https://github.com/mediocregopher/goat
[godocs]: http://godoc.org/github.com/mediocregopher/seq
[license]: /LICENSE
[examples]: /examples
