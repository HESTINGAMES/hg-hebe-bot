#!/bin/bash
set -ex

function do_fix() {
  git ls-files -- *.go | grep -v _vendor | xargs gofmt -s -w
  git ls-files -- *.go | grep -v _vendor | xargs goimports -w
}

function do_unit_test() {
  go test -race ./...
}

function do_integration_test() {
  go test -race -tags=integration ./...
}

function do_lint() {
  gometalinter --vendor --min-confidence=.3 -t --deadline=30s --disable-all -Egolint -Etest -Eineffassign -Etestify -Eunconvert -Estaticcheck -Egoconst -Egocyclo -Eerrcheck -Egofmt -Evet -Edupl -Einterfacer -Estructcheck -Evetshadow -Egosimple -Egoimports -Egolint -Egolint -Evarcheck -Emisspell -Ealigncheck -Etest  ./...
}

function do_all() {
  do_setup
  do_lint
  do_test
  do_build
}

function do_setup() {
  go get -u github.com/alecthomas/gometalinter
  gometalinter --install --update
}

case "$1" in
  fix)
    do_fix
    ;;
  test)
    do_unit_test
    ;;
  integration_test)
    do_integration_test
    ;;
  lint)
    do_lint
    ;;
  setup)
    do_setup
    ;;
  all)
    do_all
    ;;
  *)
  echo "Usage: $0 {fix|test|integration_test|lint|setup|all}"
    exit 1
    ;;
esac
