TFLIB=$(shell cd var/lib && pwd)
MOUNT_PATH=$(shell cd var && pwd)
TGO_ENV := LIBRARY_PATH=${LIBRARY_PATH}:${TFLIB} LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:${TFLIB} DYLD_LIBRARY_PATH=${DYLD_LIBRARY_PATH}:${TFLIB}
EXPORT_IMAGE := bert-export
COVERFILE := coverage.out
NUM_LABELS := 2
MODEL ?= bert-base-uncased

check: lint test

clean:
	# TODO python
	rm ${COVERFILE}

cover: cover/func

cover/%:
	go tool cover -$*=${COVERFILE}

ex/search:
	${TGO_ENV} go run ./examples/semantic-search -seqlen=16 ${MOUNT_PATH}/export/${MODEL} ./examples/semantic-search/go-faq.csv

ex/%:
	${TGO_ENV} MODEL_PATH=${MODEL_PATH} go run ./examples/$*

get:
	go get -u golang.org/x/lint/golint
	${TGO_ENV} go get ./...

image/export:
	cd export && docker build -t ${EXPORT_IMAGE} .

#inspect_model/%:
#	# TODO drop in favor of CMD
#	python ${TF_ROOT}/tensorflow/python/tools/saved_model_cli.py show --dir=$* --all

lint:
	go vet ./...
	golint ./...

model: export_image
	mkdir -p ${MOUNT_PATH}
	docker run -v ${MOUNT_PATH}:/var/bert ${EXPORT_IMAGE} export_embedding.py --download=${MODEL} /var/bert/model /var/bert/export/${MODEL}

model/classifier: export_image
	mkdir -p ${MOUNT_PATH}
	docker run -v ${MOUNT_PATH}:/var/bert ${EXPORT_IMAGE} export_classifier.py /var/bert/model/${MODEL} /var/bert/export/${MODEL} ${NUM_LABELS}

test:
	${TGO_ENV} go test -coverprofile=${COVERFILE} -v ./...
