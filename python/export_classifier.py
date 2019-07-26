# coding=utf-8

from __future__ import absolute_import
from __future__ import print_function
import sys
sys.path.insert(0,'bert')

import os.path
import argparse

import tensorflow as tf

import bert.modeling as modeling
from bert.run_classifier import create_model, flags
from util import export


parser = argparse.ArgumentParser()
parser.add_argument("model_path", help="Path for fine-tuned classifier model")
parser.add_argument("export_path", help="Path to export to")
parser.add_argument("num_labels", type=int, help="Number of labels")

parser.add_argument("--max_seq_length", type=int, default=128,
        help="Must be same as when model was fine-tuned")
parser.add_argument("--bert_config_path", help="If bert_config is not in"
        "model_path/bert_config.json, specify its path here")


def export_classifier(args):
    config_path = os.path.join(args.model_path, "bert_config.json")
    if args.bert_config_path:
        config_path = args.bert_config_path
    max_seq_length = args.max_seq_length
    num_labels = args.num_labels

    def transfer():
        bert_config = modeling.BertConfig.from_json_file(config_path)
        input_ids = tf.placeholder(tf.int32, (None, max_seq_length), name="input_ids")
        input_mask = tf.placeholder(tf.int32, (None, max_seq_length), name="input_mask")
        segment_ids = tf.placeholder(tf.int32, (None, max_seq_length), name="input_type_ids")
        # Label IDs isn't an input, so can this be moved to a constant?
        label_ids = tf.placeholder(tf.int32, (num_labels))

        _,_,_, probs = create_model(bert_config, False, input_ids, input_mask,
                                    segment_ids, label_ids, num_labels, False)
        probs = tf.identity(probs, 'probabilities')
        return {
            #'label_ids': label_ids,
            'input_ids': input_ids,
            'input_mask': input_mask,
            'input_type_ids': segment_ids
        }, {
            "probabilities": probs
        }
    export(args.model_path, args.export_path, transfer)


if __name__ == '__main__':
    args = parser.parse_args()
    export_classifier(args) 
