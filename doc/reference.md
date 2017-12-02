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

## Logic functions

(**not** val) → bool
> Returns false if the `val` is truthy (not false or nil), otherwise true.


## List functions

(**list** val<sub>1</sub> val<sub>2</sub> ... val<sub>n</sub>) → list

> Return a list constituting all the parameter values. The values can
> be any type, including lists.

(**list?** val) → bool
> Return true of the value is a list.

(**head** list) → value
> Return the first value in the list.

(**tail** list) → value
> Return the remainder of the list, ignore the first value.

(**append** list value) → list
> Add the value to the end of the list, returning a new list.

(**prepend** value list) → list
> Add the value to the beginning of the list, returning a new list.

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

## File functions

(**read-file** file-name) → string
> Return the contents of the named file as a string.

(**open** file-name) → file-handle <span style="color:red">_;; not implemented_</span>
> Open a file for reading or writing.

(**close** file-handle) → nil <span style="color:red">_;; not implemented_</span>
> Close an open file handle.

(**exists?** file-name-or-handle) → bool <span style="color:red">_;; not implemented_</span>
> Returns true if the file or directory exists.

(**file?** file-name-or-handle) => bool <span style="color:red">_;; not implemented_</span>
> Returns true of the file is a file.

(**directory?** file-name-or-handle) → bool <span style="color:red">_;; not implemented_</span>
> Returns true if the file is a directory.

(**directories** file-name-or-handle) → list <span style="color:red">_;; not implemented_</span>
> Return a list of all the files and directories (recursive) starting
at file-name-or-handle as the root.

(**read-line** file-handle) → string | nil <span style="color:red">_;; not implemented_</span>
> Read a line from a file. A `nil` signifies an end-of-file condition.
