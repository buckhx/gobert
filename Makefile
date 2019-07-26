TF_ROOT=${GOPATH}/src/github.com/tensorflow/tensorflow
TGO_ENV := LIBRARY_PATH=${TF_ROOT}/bazel-bin/tensorflow LD_LIBRARY_PATH=${TF_ROOT}/bazel-bin/tensorflow
MODEL_PATH ?= var/export/embedding_optimized

check: lint test

clean:
	# TODO flexible model
	rm -rf python/output/*
	rm coverage.out

go:
	${TGO_ENV} MODEL_PATH=${MODEL_PATH} go run main.go

ex/search:
	${TGO_ENV} go run ./examples/semantic-search var/export/embedding_optimized var/glue/QQP/original/quora.csv

ex/%:
	${TGO_ENV} MODEL_PATH=${MODEL_PATH} go run ./examples/$*

model/classifier:
	cd python && python export_embedding.py ${MODEL_PATH} var/export/classifier 2

model/embedding:
	# TODO flexible model w/ download
	cd python && python export_embedding.py ${MODEL_PATH} var/export/embedding

inspect_model/%:
	python ${TF_ROOT}/tensorflow/python/tools/saved_model_cli.py show --dir=$* --all

lint:
	go vet ./...
	golint ./...

test:
	${TGO_ENV} go test -coverprofile=coverage.out -v ./...
