FROM golang:1.12.7

# TODO buildarg
ENV TF_C_ARCHIVE=https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.14.0.tar.gz
RUN curl $TF_C_ARCHIVE | tar -C /usr/local -zxv
COPY . .
RUN go get github.com/golang/protobuf/proto && go get ./... && go build -o search .
CMD ["./search", " -seqlen=16", "/var/bert/export/bert-base-uncased", "go-faq.csv"]
