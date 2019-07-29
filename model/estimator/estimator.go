// Package estimator is a utility method for interactinfg with tf models.
// *** Experimental ***
// This package is meant ot be a pseudo-port of the python Estimator API
package estimator

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// InputFunc matches feeds in sessions
type InputFunc func(m *tf.SavedModel) map[tf.Output]*tf.Tensor

// ModelFunc the returned params match fetches & targets from the API
type ModelFunc func(m *tf.SavedModel) ([]tf.Output, []*tf.Operation)

// Estimator matches the tf, p
type Estimator interface {
	/*
		Trainer
		Evaluator
		Exporter
	*/
	Predictor
}

// Predictor creates tensors for prediction
type Predictor interface {
	Predict(InputFunc) ([]*tf.Tensor, error)
}

/*
type Evaluator interface {
	Evaluate(InputFunc) ([]*tf.Tensor, error)
}

type Trainer interface {
	Train(InputFunc) ([]*tf.Tensor, error)
}

type Exporter interface {
	Export(InputFunc) ([]*tf.Tensor, error)
}
*/
