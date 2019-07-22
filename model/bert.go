package model

import (
	"fmt"

	"github.com/buckhx/gobert/model/estimator"
	"github.com/buckhx/gobert/tokenize"
	"github.com/buckhx/gobert/vocab"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// Operation names
const (
	IDsOpName     = "input_ids"
	MaskOpName    = "input_mask"
	TypeIDsOpName = "input_type_ids"
	OutputOp      = "outputs"
)

const DefaultSeqLen = 128

/*
func NewBertClassifier(path string, opts ...BertOption) (Bert, error) {
	return NewBert(path, append(opts,
		WithModelFunc(func(m *tf.SavedModel) ([]tf.Output, []*tf.Operation) {
			return []tf.Output{
				m.Graph.Operation("SoftMax/Classifier").Output(0),
			}, nil
		}),
	)...)
}
*/

type TensorInputFunc func(map[string]*tf.Tensor) estimator.InputFunc

type ValueProvider interface {
	Value() interface{}
}

type Bert struct {
	m         *tf.SavedModel
	p         estimator.Predictor
	factory   *FeatureFactory
	modelFunc estimator.ModelFunc
	inputFunc TensorInputFunc
}

func NewBert(path string, opts ...BertOption) (Bert, error) {
	tags := []string{"bert-uncased"} // TODO configure tags
	voc, err := vocab.FromFile(path + "/vocab.txt")
	if err != nil {
		return Bert{}, err
	}
	tkz := tokenize.NewTokenizer(voc)
	fb := &FeatureFactory{tokenizer: tkz, seqLen: DefaultSeqLen}
	m, err := tf.LoadSavedModel(path, tags, nil)
	if err != nil {
		return Bert{}, err
	}
	b := Bert{
		m:       m,
		factory: fb,
		inputFunc: func(inputs map[string]*tf.Tensor) estimator.InputFunc {
			return func(m *tf.SavedModel) map[tf.Output]*tf.Tensor {
				return map[tf.Output]*tf.Tensor{
					m.Graph.Operation(IDsOpName).Output(0):     inputs[IDsOpName],
					m.Graph.Operation(MaskOpName).Output(0):    inputs[MaskOpName],
					m.Graph.Operation(TypeIDsOpName).Output(0): inputs[TypeIDsOpName],
				}
			}
		},
		modelFunc: func(m *tf.SavedModel) ([]tf.Output, []*tf.Operation) {
			return []tf.Output{
					m.Graph.Operation(OutputOp).Output(0),
				},
				nil
		},
	}
	for _, opt := range opts {
		b = opt(b)
	}
	b.p = estimator.NewPredictor(m, b.modelFunc)
	return b, nil
}

func (b Bert) PredictValues(texts ...string) ([]ValueProvider, error) {
	fs := b.factory.Features(texts...)
	inputs, err := Tensors(fs...)
	if err != nil {
		return nil, err
	}
	res, err := b.p.Predict(b.inputFunc(inputs))
	if err != nil {
		return nil, err
	}
	//	return res, nil
	vals := make([]ValueProvider, len(res))
	for i, t := range res {
		vals[i] = ValueProvider(t)
	}
	return vals, nil
}

func printModel(m *tf.SavedModel) {
	fmt.Printf("%+v\n", m)
	fmt.Println("Session")
	fmt.Println("\tDevice")
	devs, err := m.Session.ListDevices()
	if err != nil {
		fmt.Println(err)
	}
	for _, dev := range devs {
		fmt.Printf("\t\t%+v\n", dev)
	}
	fmt.Println("Graph")
	fmt.Println("\tOperation")
	for _, op := range m.Graph.Operations() {
		fmt.Printf("\t\t%+v\n", op.Name())
	}
}
