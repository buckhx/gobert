package model

import (
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

type Bert struct {
	m         *tf.SavedModel
	p         estimator.Predictor
	seqLen    int32
	tokenizer tokenize.VocabTokenizer
	modelFunc estimator.ModelFunc
	inputFunc TextInputFunc
}

func NewBert(path string, vocabPath string) (Bert, error) {
	tags := []string{"bert-uncased"} // TODO configure tags
	voc, err := vocab.FromFile(vocabPath)
	if err != nil {
		return Bert{}, err
	}
	tkz := tokenize.NewTokenizer(voc)
	m, err := tf.LoadSavedModel(path, tags, nil)
	if err != nil {
		return nil, err
	}
	b := Bert{
		m:         m,
		seqLen:    64, //TODO default
		tokenizer: tkz,
		inputFunc: func(texts ...string) estimator.InputFunc {
			fb := FeatureFactory{tokenizer: tkz, seqLen: seqLen}
			fs := fb.Features(texts...)
			inputs, err := Tensors(fs...)
			if err != nil {
				return nil, err
			}
			return func(m *tf.SavedModel) [tf.Output]*tf.Tensor {
				return map[tf.Output]*tf.Tensor{
					m.Graph.Operation(IDsOpName).Output(0):     inputs[IDsOpName],
					m.Graph.Operation(MaskOpName).Output(0):    inputs[MaskOpName],
					m.Graph.Operation(TypeIDsOpName).Output(0): inputs[TypeIDsOpName],
				}
			}
		},
		modelFunc: func(m *tf.SavedModel) ([]tf.Ouput, []*tf.Operation) {
			return []tf.Output{
					m.Graph.Operation(OutputOp).Output(0),
				},
				nil
		},
	}
	p := estimator.NewPredictor(m, b.modelFunc)
}

/*
	p := estimator.NewPredictor(m, func(m *tf.SavedModel) ([]tf.Ouput, []*tf.Operation) {
		return []tf.Output{
				b.m.Graph.Operation(OutputOp).Output(0),
			},
			nil
	})
	return Bert{
		m:         m,
		tokenizer: tkz,
		seqLen:    seqLen,
	}, nil
*/

func (b Bert) Predict(texts ...string) [][][]float32 {
	b.p.Predict(b.inputFunc(texts...))
}

type TextInputFunc func(texts ...string) estimator.InputFunc

func defaultInput(texts ...string) estimator.InputFunc {
	fb := FeatureFactory{tokenizer: b.tokenizer, seqLen: b.seqLen}
	fs := fb.Features(texts...)
	inputs, err := Tensors(fs...)
	if err != nil {
		return nil, err
	}
	return map[tf.Output]*tf.Tensor{
		m.Graph.Operation(IDsOpName).Output(0):     inputs[IDsOpName],
		m.Graph.Operation(MaskOpName).Output(0):    inputs[MaskOpName],
		m.Graph.Operation(TypeIDsOpName).Output(0): inputs[TypeIDsOpName],
	}
}

// Infer retturns an interfence from the given texts
// The text inferencecs are returned with the stame index they are supplied with
// The esecodin dimension corresponds to the token position
// TODO: decide if a wrapper API should be first-class or work w/ raw data
func (b Bert) Infer(texts ...string) ([][][]float32, error) {
	fb := FeatureFactory{tokenizer: b.tokenizer, seqLen: b.seqLen}
	fs := fb.Features(texts...)
	inputs, err := Tensors(fs...)
	if err != nil {
		return nil, err
	}
	res, err := b.m.Session.Run(
		map[tf.Output]*tf.Tensor{
			b.m.Graph.Operation(IDsOpName).Output(0):     inputs[IDsOpName],
			b.m.Graph.Operation(MaskOpName).Output(0):    inputs[MaskOpName],
			b.m.Graph.Operation(TypeIDsOpName).Output(0): inputs[TypeIDsOpName],
		},
		[]tf.Output{
			b.m.Graph.Operation(OutputOp).Output(0),
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, fmt.Errorf("Invalid Model Output Shape: [%d]", len(res))
	}
	mat, ok := res[0].Value().([][][]float32)
	if !ok {
		return nil, fmt.Errorf("Invalid Model Output Assertion to [][][]float32")
	}
	return mat, nil
	/*
		infs := make([]Inference, len(raw))
		for i, row := range raw {
			toks := make([]Vector, len(row))
			for j, v := range row {
				toks[j] = Vector(v)
			}
			infs[i] = Inference{
				tokens: toks,
			}
		}
		return infs, nil
	*/
}

func debugModel(m *tf.SavedModel) {
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
