# coding=utf-8

from __future__ import absolute_import
from __future__ import print_function
import sys
sys.path.insert(0, 'bert')  # noqa

import os.path
import argparse

import tensorflow as tf

import bert.modeling as modeling
from util import export
from download_pretrained import download


parser = argparse.ArgumentParser()
parser.add_argument("model_path", help="Path for pre-trained BERT model")
parser.add_argument("export_path", help="Path to export to")

parser.add_argument("--bert_config_path", help="If bert_config is not in"
                    "model_path/bert_config.json, specify its path here")
parser.add_argument("--download", help="Download pretrained model by name"
                    "Model be saved in model_path with name appened")


def export_embedding(args):
    # pretrained_name = os.path.basename(os.path.normpath(args.model_path))
    dl = args.download
    if dl:
        download(dl, args.model_path)
        args.model_path = os.path.join(args.model_path, dl)
        print("Model Path updated:", args.model_path)
    config_path = os.path.join(args.model_path, "bert_config.json")
    if args.bert_config_path:
        config_path = args.bert_config_path

    def transfer():
        bert_config = modeling.BertConfig.from_json_file(config_path)
        # Inputs
        # TODO shapes
        # unique_ids = tf.compat.v1.placeholder(tf.int32, (None), 'unique_ids')
        input_ids = tf.compat.v1.placeholder(tf.int32, (None, None), 'input_ids')
        input_mask = tf.compat.v1.placeholder(tf.int32, (None, None), 'input_mask')
        segment_ids = tf.compat.v1.placeholder(tf.int32, (None, None), 'input_type_ids')
        # Model
        model = modeling.BertModel(
            config=bert_config,
            is_training=False,
            input_ids=input_ids,
            input_mask=input_mask,
            token_type_ids=segment_ids,
            use_one_hot_embeddings=False)
        # Output Outputs
        # Use second to last layer, empirically performs well
        layer_ids = [-2]
        layers = tf.concat([model.all_encoder_layers[l] for l in layer_ids], -1)
        mask = tf.cast(input_mask, tf.float32)
        masker = lambda x, m: x * tf.expand_dims(m, axis=-1)
        output = masker(layers, mask)
        embedding = tf.identity(output, 'embedding')
        return {
             #  'unique_ids': unique_ids,
            'input_ids': input_ids,
            'input_mask': input_mask,
            'input_type_ids': segment_ids
        }, {
            #  "feature_ids": unique_ids,
            "embedding": embedding
        }
    export(args.model_path, args.export_path, transfer,
           method_name="bert/pretrained/embedding",
           sig_name="embedding", tags=["bert-pretrained"])


if __name__ == '__main__':
    args = parser.parse_args()
    export_embedding(args)
