<!-- mode: gfm-mode; fill-column: 78 -->
# haki

Lisp interpreter for writing scripts and hopefully as something you
could embed in Golang programs.

## status

Incomplete. Until I have a complete set of file/io functions as well
as shell-exec functions, this interpreter isn't really useful for any
real-world tasks.

## install

    $ go get -u github.com/zentrope/haki/cmd/haki

## docs

 * [function reference](doc/reference.md)

## todo

 * version info
 * ~~hash-map data structure and functions~~
 * ~~file io~~
 * ~~shell cmd exec~~
 * ~~command-line arguments~~

## non-goals

* macros
* threading
* exceptions (try/catch, call/cc, etc)

## looks

Stuff you can do at the `repl` or when invoked as a script:

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
