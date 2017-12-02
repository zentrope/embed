# haki

A rough guide to the functions you can use in the Haki scripting language.


## Math functions

(**+** num<sub>1</sub> num<sub>2</sub> ... num<sub>n</sub>) → num
> Returns the sum of all the numeric parameters.

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

(**not** val) → bool
> Returns false if the `val` is truthy (not false or nil), otherwise true.

(__=__ val<sub>1</sub> val<sub>2</sub>) → bool
> Returns true if the `val` of each of the params is equivalent
(regardless of whether or not they're the same object in memory).



## List functions

(**list** val<sub>1</sub> val<sub>2</sub> ... val<sub>n</sub>) → list

> Return a list constituting all the parameter `vals`. The `vals` can
> be any type, including lists.

(**list?** val) → bool
> Return true of the `val` is a list.

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

(**reduce** function initial-value list) → value
> Returns a value calculated by `function` applied to the
`initial-value` and each item in the `list`, with the `function` returning
a new `initial-value` per `list` item.

(**range** num) → list
> Returns a list of numbers from 0 to `num` step 1.

(**take** num list) → list
> Returns a new list consisting of the first `num` values in `list`.

(**count** list) → num
> Returns the number of items in the `list`.

(**join** list<sub>1</sub> list<sub>2</sub> ... list<sub>n</sub>) → list

> Returns the concatenation of each parameter into a single list.

(__loop__ fn list) → nil

> Applies `fn` to each val in `list` for side-effects.

(__loop-index__ fn list) → nil

> Applies `fn` to an incrementing index and each `val` in `list` (example: `(fn (idx
> val) (prn idx val))`) for side-effects.

## Print functions

(__prn__ val<sub>1</sub> val<sub>2</sub> ... val<sub>n</sub>) → nil
> Prints the values to standard out, appending a newline.


## String functions

(__format__ pattern val<sub>1</sub> val<sub>2</sub>... val<sub>n</sub>) → string

> Formats a string based on pattern and the value parameters (a.k.a,
__sprintf__) according to the [Golang implementation][printf].

(__re-find__ regex string) → string

> Returns the first match for `regex` in `string`.

(__re-list__ regex string) → list

> Returns a list of all `regex` matches in `string`.

(__re-match__ regex string) → bool

> Returns true if the `regex` finds a match in `string`.

(__re-split__ regex string) → list

> Returns a list of strings split based on the `regex` applied to `string`.


[printf]: https://golang.org/pkg/fmt/



## File functions

(__read-file__ file-name) → string

> Return the contents of the named file as a string.

(__open__ file-name) → file-handle

> Open a file for reading or writing.

(__close__ file-handle) → nil

> Close an open file handle.

(__handle?__ file-handle) → bool

> Return true if `file-handle` is a file-handle returned by open.

(__exists?__ file-namee) → bool

> Returns true if the file or directory exists.

(__file?__ file-name) => bool

> Returns true if the file is a file (not a directory).

(__directory?__ file-name) → bool

> Returns true if the file is a directory (not a file).

(__directories__ file-name) → list <span style="color:red">_;; not implemented_</span>

> Return a list of all the files and directories (recursive) starting
at file-name-or-handle as the root.

(__read-line__ file-handle) → string | nil <span style="color:red">_;; not implemented_</span>

> Read a line from a `file-handle`. A `nil` signifies an end-of-file condition.
