# embed

A small interpreter to explore the idea of using a small language to create little grammars or parsers or evaluators or rules you can add to a Golang program as data. Something you could use if a single regular expression doesn't really work all that well for you.

I'm thinking this is a kind of string processing kind of DSL. You can pass in a giant text file, and a bit of script and it can return parts of the file as a result, or build an aggregate out of it.

## todo

This experiment is designed to be something you'd use to transform data according to rules that are best expressed as regular code so I'm not going to worry too much about file IO or socket connections. You pass in a string, you get another string (or a collection of strings) back.

* [x] ~~repl~~
* [x] ~~top level definitions~~
* [x] ~~top level functions~~
* [x] ~~do expression (special)~~
* [x] ~~let special form~~
* [x] ~~prepend -- builtin list function~~
* [x] ~~append -- builtin list function~~
* [x] ~~join -- builtin list function~~
* [x] ~~anonymous "lambda" functions~~
* [x] ~~list -- builtin list function~~
* [x] ~~load core/prelude functions (written in DSL)~~
* [ ] implement prelude: ~~map~~, ~~reduce~~, ~~range~~, filter, in the DSL itself
* [ ] cond special form
* [ ] mutation
* [ ] apply builtin
* [ ] varargs or &rest parameters
* [ ] builtin: regex matching
* [ ] builtin: regex group stuff
* [ ] builtins: string functions (replace, replace-all, concat, starts, ends, trim, index, lastindex, split, ...).
* [ ] comments
* [ ] embed API for Golang programs
* [ ] tests
* [ ] load-code and load-data (handy for interactive dev/testing)
* [ ] consider a [cps](https://stackoverflow.com/a/5986168) interpreter (recursion is a killer)

## issues

* [x] ~~Repl should read all forms before presenting prompt~~
* [ ] Def/un should always store in global env.
* [x] ~~Pressing "return" in repl should not generate EOF error~~
* [x] ~~`(map (fn (x) (+ 2 x)) (list 1 2 3))` tries to eval '1~~
* [x] ~~`(map (fn (x) (+ 2 x)) '(1 2 3))` tries to eval '1~~

## looks

Stuff you can do at the `repl` as of this writing.

``` emacs-lisp
(def a 2)
(def b 3)

(defun add (x y)
  (let (i (+ a x)
        j (+ b y))
    (+ i j))

repl> (add a b)
10

repl> (add 10 b)
18

repl> (join '(1 2 3) (append '(a b) 'c) (prepend 'x '(y z)))
(1 2 3 a b c x y z)

;; Anonymous functions
(defun addf (a f)
  (+ a (f a)))

repl> (addf 1 (fn (x) (+ x 10)))
12

;; Lexical scope
(defun mk-addr (x)
  (fn (y) (+ x y)))

repl> (def add2 (mk-addr 2))
fn<fn3 (y)>

repl> (add2 1)
3

repl> (add2 1000)
1002


;; map
(defun map (f xs)
  (if (= xs '())
    xs
    (prepend (f (head xs)) (map f (tail xs)))))

repl> (map (fn (x) (+ x 2)) '(10 20 30))
(12 22 32)

repl> (map (fn (x) (+ x 2)) (list (+ 10 10) (+ 20 20) (+ 40 40)))
(22 42 82)

```

## license

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
