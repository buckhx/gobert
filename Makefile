TFLIB=$(shell cd var/lib && pwd)
TGO_ENV := LIBRARY_PATH=${LIBRARY_PATH}:${TFLIB} LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:${TFLIB} DYLD_LIBRARY_PATH=${DYLD_LIBRARY_PATH}:${TFLIB}
MODEL_PATH ?= var/export/embedding

check: lint test

clean:
	# TODO flexible model
	rm -rf python/output/*
	rm coverage.out

get:
	${TGO_ENV} go get ./...

go:
	${TGO_ENV} MODEL_PATH=${MODEL_PATH} go run main.go

ex/search:
	${TGO_ENV} go run ./examples/semantic-search -d=\t ${MODEL_PATH} var/quotes/quotes.csv

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
