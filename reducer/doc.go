/*
Package reducer implements reducers, similar to the Reducer concept in Clojure (https://clojure.org/reference/reducers).

Unfortunately, GO's generic implementation is too limited to make this abstraction pleasant to work with.
In most cases, it is recommended to use the iterable package instead.

However, reducers simplify the implementation of some algorithms, for example the lazy quicksort implemented in this package.

*/
package reducer
