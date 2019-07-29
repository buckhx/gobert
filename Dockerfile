FROM busybox:1.31.0 AS builder

# SEE https://www.tensorflow.org/install/lang_c
ARG TF_LIB_ARCHIVE=https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.14.0.tar.gz
WORKDIR /usr/local 
RUN wget -qO- $TF_LIB_ARCHIVE | tar -zxv

FROM scratch
COPY --from=builder /usr/local/lib/libtensorflow* /usr/local/lib/
