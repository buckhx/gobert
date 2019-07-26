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
	//#	UniqueIDsOp    = "unique_ids"
	InputIDsOp     = "input_ids"
	InputMaskOp    = "input_mask"
	InputTypeIDsOp = "input_type_ids"
)

const (
	DefaultSeqLen    = 128
	DefaultVocabFile = "vocab.txt"
)

type TensorInputFunc func(map[string]*tf.Tensor) estimator.InputFunc
type FeatureTensorFunc func(fs ...tokenize.Feature) (map[string]*tf.Tensor, error)

type ValueProvider interface {
	Value() interface{}
}

type Bert struct {
	m          *tf.SavedModel
	p          estimator.Predictor
	factory    *tokenize.FeatureFactory
	modelFunc  estimator.ModelFunc
	inputFunc  TensorInputFunc
	tensorFunc FeatureTensorFunc
}

// Pipeline: text -> FeatureFactory -> TensorFunc -> InputFunc -> ModelFunc -> Value
func NewBert(m *tf.SavedModel, vocabPath string, opts ...BertOption) (Bert, error) {
	voc, err := vocab.FromFile(vocabPath)
	if err != nil {
		return Bert{}, err
	}
	tkz := tokenize.NewTokenizer(voc)
	b := Bert{
		m:          m,
		factory:    &tokenize.FeatureFactory{Tokenizer: tkz, SeqLen: DefaultSeqLen},
		tensorFunc: tensors,
		inputFunc: func(inputs map[string]*tf.Tensor) estimator.InputFunc {
			return func(m *tf.SavedModel) map[tf.Output]*tf.Tensor {
				return map[tf.Output]*tf.Tensor{
					//	m.Graph.Operation(UniqueIDsOp).Output(0):    inputs[UniqueIDsOp],
					m.Graph.Operation(InputIDsOp).Output(0):     inputs[InputIDsOp],
					m.Graph.Operation(InputMaskOp).Output(0):    inputs[InputMaskOp],
					m.Graph.Operation(InputTypeIDsOp).Output(0): inputs[InputTypeIDsOp],
				}
			}
		},
		modelFunc: func(m *tf.SavedModel) ([]tf.Output, []*tf.Operation) {
			return []tf.Output{
					m.Graph.Operation(EmbeddingOp).Output(0),
					//		m.Graph.Operation("feature_ids").Output(0),
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

func (b Bert) Features(texts ...string) []tokenize.Feature {
	return b.factory.Features(texts...)
}

func (b Bert) PredictValues(texts ...string) ([]ValueProvider, error) {
	fmt.Println("Building Features...")
	fs := b.factory.Features(texts...)
	inputs, err := b.tensorFunc(fs...)
	if err != nil {
		return nil, err
	}
	fmt.Println("Done Building")
	fmt.Println("Predicting...")
	res, err := b.p.Predict(b.inputFunc(inputs))
	if err != nil {
		return nil, err
	}
	fmt.Println("Done TF Predicting")
	//	return res, nil
	vals := make([]ValueProvider, len(res))
	for i, t := range res {
		vals[i] = ValueProvider(t)
	}
	fmt.Println("Done Value Casting")
	return vals, nil
}

func Print(m *tf.SavedModel) {
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
		fmt.Printf("\t\t%s %s\t%d/%d\n", op.Name(), op.Type(), op.NumInputs(), op.NumOutputs())
		for i := 0; i < op.NumOutputs(); i++ {
			o := op.Output(i)
			fmt.Printf("\t\t\t%v %s\n", o.DataType(), o.Shape())
		}
	}
}

func tensors(fs ...tokenize.Feature) (map[string]*tf.Tensor, error) {
	//	uids := make([]int32, len(fs))
	tids := make([][]int32, len(fs))
	mask := make([][]int32, len(fs))
	sids := make([][]int32, len(fs))
	for i, f := range fs {
		//		uids[i] = f.ID
		tids[i] = f.TokenIDs
		mask[i] = f.Mask
		sids[i] = f.TypeIDs
	}
	/*
		u, err := tf.NewTensor(uids)
		if err != nil {
			return nil, err
		}
	*/
	t, err := tf.NewTensor(tids)
	if err != nil {
		return nil, err
	}
	m, err := tf.NewTensor(mask)
	if err != nil {
		return nil, err
	}
	s, err := tf.NewTensor(sids)
	if err != nil {
		return nil, err
	}
	return map[string]*tf.Tensor{
		//UniqueIDsOp:    u,
		InputIDsOp:     t,
		InputMaskOp:    m,
		InputTypeIDsOp: s,
	}, nil
}
