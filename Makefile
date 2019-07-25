TF_ROOT=${GOPATH}/src/github.com/tensorflow/tensorflow
TGO_ENV := LIBRARY_PATH=${TF_ROOT}/bazel-bin/tensorflow LD_LIBRARY_PATH=${TF_ROOT}/bazel-bin/tensorflow
GOBERT_BASE_DIR ?= var/export/embedding_optimized

check: lint test

clean:
	# TODO flexible model
	rm -rf python/output/*
	rm coverage.out

go:
	${TGO_ENV} GOBERT_BASE_DIR=${GOBERT_BASE_DIR} go run main.go

ex/embedding:
	${TGO_ENV} GOBERT_BASE_DIR=${GOBERT_BASE_DIR} go run examples/embedding/main.go

search:
	${TGO_ENV} go run ./examples/semantic-search var/export/embedding_optimized var/glue/QQP/original/quora.csv

ex/%:
	${TGO_ENV} GOBERT_BASE_DIR=${GOBERT_BASE_DIR} go run examples/$*/main.go

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
