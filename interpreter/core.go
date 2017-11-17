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
	(if (= xs '())
		xs
		(let (x (head xs))
			(if (f x)
				(prepend x (filter f (tail xs)))
				(filter f (tail xs))))))

(defun range (x)
	(if (= x 0) (list 0)
		(append (range (- x 1)) x)))

(defun even? (x)
	(= (mod x 2) 0))

(defun odd? (x)
	(not (even? x)))`
