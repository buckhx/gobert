TF_ROOT=${GOPATH}/src/github.com/tensorflow/tensorflow
TGO_ENV := LIBRARY_PATH=${TF_ROOT}/bazel-bin/tensorflow LD_LIBRARY_PATH=${TF_ROOT}/bazel-bin/tensorflow

go:
	${TGO_ENV} go run main.go

inspect_model/%:
	python ${TF_ROOT}/tensorflow/python/tools/saved_model_cli.py show --dir=$* --all
