# gobert

Go bindings for operationaling BERT models. Train in Python, run in Go.

The following advice from Go TensorFlow applies:
```
TensorFlow provides a Go API— particularly useful for loading models created with Python and running them within a Go application.
...
TensorFlow provides APIs for use in Go programs. These APIs are particularly well-suited to loading models created in Python and executing them within a Go application.
...
Be a real gopher, keep it simple! Use Python to define & train models; you can always load trained models and using them with Go later!
```

# TBD
[ ] gonum interop
[ ] first class wrapper API ([][][]float32 -> []Inference)
[ ] proto interops
[ ] pooling strategies

Current line of thought is to use core lib for raw types and supply a utility API

# Models

From pytho

[ ] BertModel
[ ] BertForPreTraining
[ ] BertForMaskedLM
[ ] BertForNextSentencePrediction
[ ] BertForSequenceClassification
[ ] BertForMultipleChoice
[ ] BertForTokenClassification
[ ] BertForQuestionAnswering
