# set variables
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
tool_name := $(current_dir)
build_dir := bin
go_bin := go
full_binary_path := $(dir $(mkfile_path))$(build_dir)/$(tool_name)
coverage_path=./tmp/go-test-coverage.out

GIT_COMMIT=$$(git rev-list -1 HEAD)
GIT_BRANCH=$$(git rev-parse --abbrev-ref HEAD)
GIT_VERSION=$$(git rev-list --count HEAD)
GIT_DATE=$$(git show -s --format=%ci HEAD)
GIT_STATE=$$(git diff --quiet && echo 'clean' || echo 'dirty')
GIT_REMOTE=$$(git config --get remote.origin.url)

.PHONY : build
build:
	$(go_bin) generate ./...
	# due to sqlite
	CGO_ENABLED=1 $(go_bin) build -asmflags -trimpath -ldflags "-buildid= -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GitBranch=$(GIT_BRANCH)' -X 'main.GitDate=$(GIT_DATE)' -X 'main.GitVersion=$(GIT_VERSION)' -X 'main.GitState=$(GIT_STATE)'  -X 'main.GitRemote=$(GIT_REMOTE)'" -o "$(full_binary_path)"

.PHONY : install
install:
	$(go_bin) generate ./...
	# due to sqlite
	CGO_ENABLED=1 $(go_bin) install -asmflags -trimpath -ldflags "-buildid= -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GitBranch=$(GIT_BRANCH)' -X 'main.GitDate=$(GIT_DATE)' -X 'main.GitVersion=$(GIT_VERSION)' -X 'main.GitState=$(GIT_STATE)'  -X 'main.GitRemote=$(GIT_REMOTE)'"

.PHONY : minify
minify:
	upx -5 "$(full_binary_path)"

.PHONY : clean
clean:
	$(go_bin) clean
	rm --force "$(full_binary_path)"
	rm --force "$(coverage_path)"
	rm --force nohup.out

.PHONY : test
test:
	$(go_bin) test

.PHONY : lint
lint:
	# https://github.com/golangci/golangci-lint
	golangci-lint run --enable-all

.PHONY : check
check: test lint

.PHONY : cover
cover:
	$(go_bin) test -coverprofile "$(coverage_path)"
	$(go_bin) tool cover -html="$(coverage_path)"
