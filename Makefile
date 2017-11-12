##
## Copyright (C) 2017 Keith Irwin
##
## This program is free software: you can redistribute it and/or modify
## it under the terms of the GNU General Public License as published
## by the Free Software Foundation, either version 3 of the License,
## or (at your option) any later version.
##
## This program is distributed in the hope that it will be useful,
## but WITHOUT ANY WARRANTY; without even the implied warranty of
## MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
## GNU General Public License for more details.
##
## You should have received a copy of the GNU General Public License
## along with this program.  If not, see <http://www.gnu.org/licenses/>.

BINARY = embed
TREE = tree

.PHONY: help init vendor
.DEFAULT_GOAL := help

godep:
	@hash dep > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -v -u github.com/golang/dep/cmd/dep; \
	fi

vendor: godep ## Make sure vendor dependencies are present.
	dep ensure

build: vendor ## Build an executable binary.
	go build -o $(BINARY)

clean: ## Clean build artifacts (vendor left alone).
	rm -f $(BINARY)

dist-clean: clean ## Clean everything, including vendor.
	rm -rf vendor

tree: ## View source hierarchy without vendor pkgs
	$(TREE) -C -I "vendor"

help: ## Produce this list of goals
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' | \
		sort
