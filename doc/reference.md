# haki

A rough guide to the functions you can use in the Haki scripting
language. The basic goals for the langage are:

* Be able to write simple, throw-away shell-script-like utilities.

* Remove the scripting (file i/o, etc) stuff to allow for embedding
  the language in `Go` programs.

The features of this language are way more than I generally need, but
way less than is generally needed.


## Special values

__nil__

> Represents the absence-of-value value, or an empty list, or non-truth.

__true__

> Represents a boolean `true` value

__false__

> Represents a boolean `false` value.

__&stdin__

> Represents the `stdin` file-handle.

__&stdout__

> Represents the `stdout` file-handle.

__&stderr__

> Represents the `stderr` file-handle.

## Math functions

(**+** num<sub>1</sub> num<sub>2</sub> ... num<sub>n</sub>) → num
> Returns the sum if all the numeric parameters.

(**-** num<sub>1</sub> num<sub>2</sub> ... num<sub>n</sub>) → num
> Returns the result of subtracting each number from the previous number.


(__*__ num<sub>1</sub> num<sub>2</sub> ... num<sub>n</sub>) → num
> Returns the product of multiplying each parameter from left to right.

(**mod** num div) → num
> Modulus of num and div.

(**inc** num) -> num
> Returns the `num` incremented by 1.

(**dec** num) → num
> Returns `num` decremented by 1.

(**even?** num) -> bool
> Returns true if `num` is an even number.

(**odd?** num) → bool
> Returns true if `num` is an odd number.


(**<** num<sub>1</sub> num<sub>2</sub> ... num<sub>n</sub>) → bool
> Returns true if each `num` is less than the number to its right.


## Logic functions

(__=__ val<sub>1</sub> val<sub>2</sub>) → bool

> Returns true if the `val` of each of the params is equivalent
(regardless of whether or not they're the same object in memory).

(__false?__ val) → bool

> Returns true if the val evaluates to a literal false.

(__nil?__ val) → bool

> Returns true if `val` is nil.

(__not__ val) → bool

> Returns false if the `val` is truthy (not false or nil), otherwise
> true.

(__true?__ val) → bool

> Returns true if the val evaluates to a literal true (not just
> truthy).

## List functions

(**list** val<sub>1</sub> val<sub>2</sub> ... val<sub>n</sub>) → list

> Return a list constituting all the parameter `vals`. The `vals` can
> be any type, including lists.

(**list?** val) → bool
> Return true if the `val` is a list.

(**head** list) → val
> Return the first value in the list.

(**tail** list) → val
> Return the remainder of the list, ignore the first value.

(**append** list val) → list
> Add the `val` to the end of the list, returning a new list.

(**prepend** val list) → list
> Add the `val` to the beginning of the list, returning a new list.

(**map** function list) → list
> Returns a list made up of `function` applied to each value in `list`.

(**filter** function list) → list
> Returns a list made up of all the values in `list` for which
(`function` item) returns true.

(**reduce** function initial-value list) → val
> Returns a value calculated by `function` applied to the
`initial-value` and each item in the `list`, with the `function` returning
a new `initial-value` per `list` item.

(**range** num) → list
> Returns a list of numbers from 0 to `num` step 1.

(**take** num list) → list
> Returns a new list consisting of the first `num` values in `list`.

(__count__ list) → int

> Returns the number of elements in the `list`.

(**join** list<sub>1</sub> list<sub>2</sub> ... list<sub>n</sub>) → list

> Returns the concatenation of each parameter into a single list.

(__loop__ fn list) → nil

> Applies `fn` to each val in `list` for side-effects.

(__loop-index__ fn list) → nil

> Applies `fn` to an incrementing index and each `val` in `list` (example: `(fn (idx
> val) (prn idx val))`) for side-effects.

## Hash Map Functions

(__count__ hash-map) → int

> Returns the number of key/value pairs in the `hash-map`.

(__hmap__ k<sub>1</sub> v<sub>1</sub> ... k<sub>n</sub> v<sub>n</sub>) → hash-map

> Construct a hash-map based on the list of `k`s and `v`s.

(__hmap?__ val) → bool

> Return true if `val` is of type `hash-map`.

(__hkeys__ m) → list

> Return all the keys in the hash-map `m` as a list.

(__hvals__ m) → list

> Return all the vals in the hash-map `m` as a list.

(__hget__ m k) → val _or_ nil

> Return the val found at position `k` in the `hash-map` or `nil` if
> not found.

(__hget-in__ m '(k<sub>1</sub> ... k<sub>n</sub>)) → val _or_ nil

> Return the val found at each value of `k` following a path through a
> map of maps, or `nil` if not found.

(__hset__ m k<sub>1</sub> v<sub>1</sub> ... k<sub>n</sub> v<sub>n</sub>) → hash-map

> Return a new hash-map adding each `k/v` pair to the old
> `hash-map`. Setting `k` to `nil` deletes the map entry.

(__hset-in__ m '(k<sub>1</sub> ... k<sub>n</sub>) v) → hash-map

> Set value of k<sub>n</sub> to val `v`, creating intermediate paths
> from k<sub>1</sub> as needed. It's an error if one of the path
> elements is present and not a hash-map.

## Print functions

(__prn__ val<sub>1</sub> val<sub>2</sub> ... val<sub>n</sub>) → nil
> Prints the values to standard out, appending a newline.


## String functions

Note: Whitespace in the following is defined as: `[' ', '\n', '\r', '\t']`.

(__count__ string) → int

> Returns the number of characters in the `string`.

(__ends-with?__ string suffix) → bool

> Returns true if string ends with suffix.

(__format__ pattern val<sub>1</sub> val<sub>2</sub>... val<sub>n</sub>) → string

> Formats a string based on pattern and the value parameters (a.k.a,
__sprintf__) according to the [Golang implementation][printf].

[printf]: https://golang.org/pkg/fmt/

(__index__ string substr) → int

> Returns the index of the first instance of `substr` in `s`, or -1
> if `substr` is not present in `s`.

(__last-index__ string substr) → int

> Returns the index of the last instance of `substr` in `s`, or -1 if
> `substr` is not present in `s`.

(__lower-case__ string) → string

> Return a new string with all letters in lower case.

(__re-find__ regex string) → string

> Returns the first match for `regex` in `string`.

(__re-list__ regex string) → list

> Returns a list of all `regex` matches in `string`.

(__re-match__ regex string) → bool

> Returns true if the `regex` finds a match in `string`.

(__re-split__ regex string) → list

> Returns a list of strings split based on the `regex` applied to
> `string`.

(__replace__ string old new) → string

> Return a copy of `string` with every instance of `old` replaced by
> `new`.

(__substr__ string start end) → string

> Return the substring of `string` starting a index `start` and ending
> at `end` (exclusive).

(__trim__ string) → string

> Trim whitespace from both ends of a string, returning a new string.

(__triml__ string) → string

> Trim whitespace from beginning of a string, returning a new string.

(__trimr__ string) → string

> Trim whitespace from the end of a string, returning a new string.

(__starts-with?__ string prefix) → string

> Returns true if string `starts` with `prefix`.

(__upper-case__ string) → string

> Return a new string with all letters in upper case.

(__words__ string) → list

> Return a list of words split from `string` using whitespace
> delimiters.

## File functions

(__close!__ file-handle) → nil

> Close an open file handle.

(__closed?__ file-handle) → bool

> Return true if the `file-handle` is closed.

(__dir?__ file-name) → bool

> Returns true if the file is a directory (not a file).

(__exists?__ file-namee) → bool

> Returns true if the file or directory exists.

(__file?__ file-name) => bool

> Returns true if the file is a file (not a directory).

(__files__ path [glob]) → list

> Return a list of all the files and directories (recursive) starting
at path. If `glob` is provided, results are filtered by matching file
names. For example: `(files "/usr/local/Cellar" "INSTALL*json")`.

(__handle?__ file-handle) → bool

> Return true if `file-handle` is a file-handle returned by open.

(__open!__ file-name) → file-handle

> Open a file for reading or writing.

(__read-file__ file-name) → string

> Return the contents of the named file as a string.

(__read-line__ file-handle) → string \ nil

> Read a line from a `file-handle`. A `nil` signifies an end-of-file
> condition.
