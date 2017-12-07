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

package lang

import "regexp"

var respace = regexp.MustCompile("\t")

func spacify(s string) string {
	return respace.ReplaceAllString(s, "  ")
}

// Core functions
var Core = spacify(`
(defun map (f xs)
	(if (= xs '())
		'()
		(prepend (f (head xs)) (map f (tail xs)))))

(defun reduce (f a xs)
	(if (= xs '())
		a
		(reduce f (f a (head xs)) (tail xs))))

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

(defun range (x)
	(let (_range (fn (x)
								 (if (= x 0)
									 (list 0)
									 (append (_range (- x 1)) x))))
		(_range (- x 1))))

(defun factorial (n)
	(let (_fact (fn (product n)
								(if (< n 2)
									product
									(_fact (* product n) (- n 1)))))
		(_fact 1 n)))

(defun take (x lst)
	(let (_take (fn (accum ls)
								(if (or (= ls '()) (= (count accum) x))
									accum
									(_take (append accum (head ls)) (tail ls)))))
		(_take '() lst)))

(defun even? (x)
	(= (mod x 2) 0))

(defun odd? (x)
	(not (even? x)))

(defun loop (f lst)
	(let (_loop (fn (f lst)
									 (if (= lst '())
										 nil
										 (do (f (head lst))
												 (_loop f (tail lst))))))
		(_loop f lst)))

(defun loop-index (f lst)
	(let (_loop (fn (index f lst)
										(if (= lst '())
											nil
											(do (f index (head lst))
													(_loop (inc index) f (tail lst))))))
		(_loop 0 f lst)))` + CoreStringFunctions)
