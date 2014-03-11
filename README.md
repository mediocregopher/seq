# seq

Clojure-like immutable data-structures and lazy lists for go.

## Disclaimer

This is not really intended for daily use. It's more of a toy and possibly
something that could be built-upon to make something cooler. Go wasn't really
built with generics in mind (and for good reason), so you lose a lot of safety
and add a lot of ugliness when you use them.

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

Check the [godocs][godocs] for the full API. Also, check out the Examples
section below.

[Examples][examples] are a good place to look too.

## About

This library constitutes an attempt at bringing immutability and laziness to go
in a thread-safe way.

### Immutability

If you have any Seq data-structure, be it a `List`, `Set`, or `HashMap`, and you
perform an operation on it (say adding an element to it) a new instance will be
returned to you with the change, leaving the original in-place. Conceptually
it's as if a copy is made on every operation. In reality Seq uses
structure-sharing to minimize the amount of copying needed. For `Set`s and
`HashMap`s there should be only two or three node copies per operation.
Additionally, the actual values being held aren't being copied.

There are multiple advantages to immutability, and most of them have been
described in great depth elsewhere. Primary benefits are:

* Code becomes easier to think about. You never have to worry about whether or
  not a function you passed a variable into changed that variable. On the flip
  side, when inside a function you never have to worry about whether or not
  you're allowed to modify a variable.

* Thread-safety comes free. Pass around variables with great abandon, there
  won't ever be race-conditions caused by two threads modifying the same value
  at the same time.

### Laziness

Seq provides a common interface for all of its structures so that they can all
be treated as sequential lists of elements. With this, Seq also provides
multiple functions for iterating and modifying these sequences, such as `Map`,
`Reduce`, `Filter`, and so on. These correspond to similar functions in other
object oriented languages.

Where possible Seq also provides lazy forms of these functions (`LMap`,
`LFilter`, etc...). These functions return a special `Lazy` type, which also
implements the `Seq` interface. With `Lazy`, the next item in the sequence won't
be evalutated until it's actually needed. So if you have a lazy map over a list
but you only consume half the result, the map function will only be called on
half the elements.

Lazy sequences cache their results in a thread-safe way. So if you pass
a reference to the same `Lazy` to multiple threads and they all consume it, each
element is only evalutated once globally.

# Legal

The [license][license] supplied in the repo applies to all code and resources
provided in the repo.

Copyright Brian Picciano, 2014

[goat]: https://github.com/mediocregopher/goat
[godocs]: http://godoc.org/github.com/mediocregopher/seq
[license]: /LICENSE
[examples]: /examples
