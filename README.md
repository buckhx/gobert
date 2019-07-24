# gobert
![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg)
[![GoDoc](https://godoc.org/github.com/buckhx/gobert?status.svg)](https://godoc.org/github.com/buckhx/gobert)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go bindings for operationalizing BERT models. Train in Python, run in Go.

The following advice from Go TensorFlow applies:
```
TensorFlow provides a Go API— particularly useful for loading models created with Python and running them within a Go application.
...
TensorFlow provides APIs for use in Go programs. These APIs are particularly well-suited to loading models created in Python and executing them within a Go application.
...
Be a real gopher, keep it simple! Use Python to define & train models; you can always load trained models and using them with Go later!
```

#
[X] Python Embedding
[X] Python Classifier
[ ] Go Classifier
[ ] Raw Model Example
[ ] Token Lookup
[ ] 
[ ] Model Download
[ ] Documentation
[ ] Test Coverage 
[ ] Benchmark
[ ] Binary CMD



# TBD
[ ] gonum interop
[ ] first class wrapper API ([][][]float32 -> []Inference)
[ ] proto interops
[ ] pooling strategies
[ ] batching

Current line of thought is to use core lib for raw types and supply a utility API

# Models

From pytho

[ ] BertModel
[ ] BertForPreTraining
[ ] BertForMaskedLM
[ ] BertForNextSentencePrediction
[X] BertForSequenceClassification
[ ] BertForMultipleChoice
[ ] BertForTokenClassification
[ ] BertForQuestionAnswering
