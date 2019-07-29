// Package estimator is a utility method for interactinfg with tf models
// This package is meant ot be a pseudo-port of the python Estimator API
package estimator

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type predictor struct {
	m       *tf.SavedModel
	outputs []tf.Output
	targets []*tf.Operation
}

// NewPredictor creates a new Predictor in lieu of a full estimator
func NewPredictor(m *tf.SavedModel, fn ModelFunc) Predictor {
	outputs, targets := fn(m)
	return predictor{
		m:       m,
		outputs: outputs,
		targets: targets,
	}
}

// Predictor will apply fn to the estimator model
func (p predictor) Predict(fn InputFunc) ([]*tf.Tensor, error) {
	inputs := fn(p.m)
	return p.m.Session.Run(inputs, p.outputs, p.targets)
}
