# gobert
![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg)
[![GoDoc](https://godoc.org/github.com/buckhx/gobert?status.svg)](https://godoc.org/github.com/buckhx/gobert)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go bindings for operationalizing [BERT](https://github.com/google-research/bert) models. Train in Python, run in Go.

Simply put, gobert translates text sentences from any language into fixed length vectors called "embeddings".
These embeddings can be used for downstream learning tasks or directly for comparison.

### BERT

BERT is a state of the art NLP model that can be leveraged for transfer learning into domain specific use cases.

* [BERT Paper](https://arxiv.org/abs/1810.04805)
* [Google AI Blog Post](https://ai.googleblog.com/2018/11/open-sourcing-bert-state-of-art-pre.html)
* [Illustrated BERT](http://jalammar.github.io/illustrated-bert/)

# Under Active Development

This is a work in progress and should not used until a version has be tagged and a go.mod is present.
Test coverage will also be added when the API settles.

The following advice from Go TensorFlow applies:
```
TensorFlow provides a Go API— particularly useful for loading models created with Python and running them within a Go application.
...
TensorFlow provides APIs for use in Go programs. These APIs are particularly well-suited to loading models created in Python and executing them within a Go application.
...
Be a real gopher, keep it simple! Use Python to define & train models; you can always load trained models and using them with Go later!
```

This project attempts to minimize dependencies

# Installation

## Prereqs

1. [Install Tensorflow for C](https://www.tensorflow.org/install/lang_c)
2. Install Docker (Optional, but suggested)

## Run Demo

The following demo will run a simple semantic search engine against the Go FAQ.
```
# Download & Export Pre-Trained Model
make model

# Run semantic search examples
make ex/search
```

### Notes

* SeqLen has large impact on performance
* Perf is not great, need to determine if it's from python model or go runtime
* model package is WIP

## Examples

* [SemanticSearch](examples/semantic-search): Simple search engine from CSV data using BERT sentence vectors
* [Classifier](examples/classifier/main.go): exposing model from run_classifier
* [Embedding](examples/embedding/main.go): returning sentence embeddings
* [Raw](examples/raw-model/main.go): Using only the gobert tokenize package and vanilla tensorflow API

## Packages

### Tokenize

The tokenize package includes methods to create BERT input features. It is fairly stable and can be used independently of the model package.
This will be its own module since it does not require tensorflow bindings.

#### Vocab

The vocab package is a simple container for BERT vocabs. Could be rolled into tokenize.

###  Model

The model package is an experimental package to work with models exported. Requires tensorflow.

There are two main external components that are required to leverage the model package. Utilities to interop with these are supplied with in this repo.

* Tensflow C Lib
* TF Model exported with the SavedModel API

### Export

The export dir includes utilities to export BERT models that can be exposed to the GO runtime.
There is a loose coupling with the model package and exported models interop with the model package.
The suggested way to run exports is through a container with a host mounted volume.

The models exported using this package interop with tensorflow/serving


# TODOs
- [X] Python Embedding
- [X] Python Classifier
- [ ] Go Classifier
- [X] Semantic Search Example
- [X] Raw Model Example
- [ ] Token Lookup
- [X] Model Download
- [ ] Documentation
- [X] Cleanup makefile
- [ ] Test Coverage
- [ ] Benchmark
- [ ] Binary CMD
- [X] Docker Exporter
- [ ] Docker TF-GO Image
- [ ] go mod init

# TBD
- [ ] Pool layers in python or post-process
- [ ] gonum interop
- [ ] first class wrapper API ([][][]float32 -> []Embedding)
- [ ] proto interops
- [ ] pooling strategies
- [ ] batching
- [ ] other tuned models (SentencePrediction, SQUAD)

Current line of thought is to use core lib for raw types and supply a utility API
