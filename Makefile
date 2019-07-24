TF_ROOT=${GOPATH}/src/github.com/tensorflow/tensorflow
TGO_ENV := LIBRARY_PATH=${TF_ROOT}/bazel-bin/tensorflow LD_LIBRARY_PATH=${TF_ROOT}/bazel-bin/tensorflow
GOBERT_BASE_DIR ?= /tmp/model/bert/export/embedding

check: lint test

clean:
	# TODO flexible model
	rm -rf python/output/*
	rm coverage.out

go:
	${TGO_ENV} GOBERT_BASE_DIR=${GOBERT_BASE_DIR} go run main.go

ex/embedding:
	${TGO_ENV} GOBERT_BASE_DIR=${GOBERT_BASE_DIR} go run examples/embedding/main.go

ex/classifier:
	${TGO_ENV} GOBERT_BASE_DIR=${GOBERT_BASE_DIR} go run examples/classifier/main.go

.PHONY: model
model:
	# TODO flexible model w/ download
	cd python && python export.py

inspect_model/%:
	python ${TF_ROOT}/tensorflow/python/tools/saved_model_cli.py show --dir=$* --all

lint:
	go vet ./...
	golint ./...

test:
	${TGO_ENV} go test -coverprofile=coverage.out -v ./...
