# embed

A small interpreter to explore the idea of using a small language to create little grammars or parsers or evaluators or rules you can add to a Golang program as data. Something you could use if a single regular expression doesn't really work all that well for you.

I'm thinking this is a kind of string processing kind of DSL. You can pass in a giant text file, and a bit of script and it can return parts of the file as a result, or build an aggregate out of it.

## todo

This experiment is designed to be something you'd use to transform data according to rules that are best expressed as regular code so I'm not going to worry too much about file IO or socket connections. You pass in a string, you get another string (or a collection of strings) back.

* [x] ~~repl~~
* [x] ~~top level definitions~~
* [x] ~~top level functions~~
* [x] ~~do expression (special)~~
* [ ] let expression (special)
* [ ] anonymous "lambda" functions
* [ ] prelude: map, reduce, filter, etc, written in the DSL itself
* [ ] mutation
* [ ] apply primitive
* [ ] varargs or &rest parameters
* [ ] primitive: regex matching
* [ ] primitive: regex group stuff
* [ ] primitives: string functions (replace, replace-all, concat, starts, ends, trim, index, lastindex, ...).
* [ ] comments
* [ ] embed API for Golang programs
* [ ] tests
* [ ] load-code and load-data (handy for interactive dev/testing)
* [ ] Fix: repl should read all forms before presenting prompt

## looks

Stuff you can do at the `repl` as of this writing.

``` emacs-lisp
(def a 2)
(def b 3)

(defun add (x y)
  (+ x y))

=> (add a b)
5

=> (add 10 b)
13
```

## License

Copyright (c) 2017 Keith Irwin

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published
by the Free Software Foundation, either version 3 of the License,
or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see
[http://www.gnu.org/licenses/](http://www.gnu.org/licenses/).
