// Package estimator is a utility method for interactinfg with tf models
// This package is meant ot be a pseudo-port of the python Estimator API
package estimator

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// Input, matches feeds in sessions
type InputFunc func(m *tf.SavedModel) map[tf.Output]*tf.Tensor

// odelFunc, the returned params match fetches & targers from the API
type ModelFunc func(m *tf.SavedModel) ([]tf.Output, []*tf.Operation)

type Estimator interface {
	Trainer
	Evaluator
	Predictor
	Exporter
}

type Predictor interface {
	Predict(InputFunc) ([]*tf.Tensor, error)
}

type Evaluator interface {
	Evaluate(InputFunc) ([]*tf.Tensor, error)
}

type Trainer interface {
	Train(InputFunc) ([]*tf.Tensor, error)
}

type Exporter interface {
	Export(InputFunc) ([]*tf.Tensor, error)
}
