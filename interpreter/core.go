//
// Copyright © 2017-present Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package interpreter

// Core functions
const Core = `
(defun map (f xs)
	(if (= xs '())
		xs
		(prepend (f (head xs)) (map f (tail xs)))))

(defun reduce (f a xs)
	(if (= xs '())
		a
		(reduce f (f a (head xs)) (tail xs))))

(defun filter (f xs)
	(prn "-------->" xs)
	(if (= xs '())
		xs
		(let (x (head xs)
					y (tail xs))
			(prn "xxxxxxxxxxxxxxxxxxxxxxxxxx " x)
			(prn "yyyyyyyyyyyyyyyyyyyyyyyyyy " y)
			(if (f x)
				(prepend x (filter f y))
				(filter f y)))))

(defun filter (f xs)
	(if (= xs '())
		xs
		(let (x (head xs)
					y (tail xs))
			(if (f (head xs))
				(prepend (head xs) (filter f (tail xs)))
				(filter f (tail xs))))))

(defun dec (x)
	(- x 1))

(defun inc (x)
	(+ x 1))

(defun range' (x)
	(if (= x 0)
		(list 0)
		(append (range' (- x 1)) x)))

(defun range (x)
	(range' (- x 1)))

(defun factorial (n)
	(factorial' 1 n))

(defun factorial' (product n)
	(if (< n 2)
		product
		(factorial' (* product n) (- n 1))))

(defun take (x lst)
	(let (take2 (fn (accum ls)
								(if (or (= ls '()) (= (count accum) x))
									accum
									(take2 (append accum (head ls)) (tail ls)))))
		(take2 '() lst)))

(defun even? (x)
	(= (mod x 2) 0))

(defun odd? (x)
	(not (even? x)))`

// This won't work until "let" allows for recursive definitions.
const deck = `
(defun take (x lst)
	(let (take' (fn (a ls)
								(if (= ls '())
										a
										(if (= (count ls) x)
											a
											(take' (append a (head ls)) (tail ls)))))
		(take' '() lst)))
`
