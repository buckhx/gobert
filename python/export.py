from __future__ import print_function

import os
import tensorflow as tf
import bert.modeling as modeling
from tensorflow.python.tools.optimize_for_inference_lib import optimize_for_inference as optimize_graph


path = os.getenv("BERT_BASE_DIR")
init_checkpoint = path+"/bert_model.ckpt"
bert_config_file = path+"/bert_config.json"
out_dir = "output"
tag = "bert-uncased"
layer_ids = [-2]

builder = tf.compat.v1.saved_model.Builder(out_dir)

# Init graph
# Config
bert_config = modeling.BertConfig.from_json_file(bert_config_file)
# Inputs
# TODO shapes
input_ids = tf.compat.v1.placeholder(tf.int32, (None, None), 'input_ids')
input_mask = tf.compat.v1.placeholder(tf.int32, (None, None), 'input_mask')
input_type_ids = tf.compat.v1.placeholder(tf.int32, (None, None), 'input_type_ids')
inputs = [input_ids, input_mask, input_type_ids]
# Model
model = modeling.BertModel(
    config=bert_config,
    is_training=False,
    input_ids=input_ids,
    input_mask=input_mask,
    token_type_ids=input_type_ids,
    use_one_hot_embeddings=False)
# Output Outputs
layers = tf.concat([model.all_encoder_layers[l] for l in layer_ids], -1)
mask = tf.cast(input_mask, tf.float32)
masker = lambda x, m: x * tf.expand_dims(m, axis=-1)
output = masker(layers, mask)
outputs = [tf.identity(output, 'outputs')]
# Restore
tvars = tf.compat.v1.trainable_variables()
amap, avars = modeling.get_assignment_map_from_checkpoint(tvars, init_checkpoint)
tf.compat.v1.train.init_from_checkpoint(init_checkpoint, amap)
# Optimize
with tf.Session() as sess:
    # TODO slim down grpah
    sess.run(tf.compat.v1.global_variables_initializer())
    g = tf.compat.v1.get_default_graph().as_graph_def()
    g = optimize_graph(g,
            [t.name[:-2] for t in inputs],
            [t.name[:-2] for t in outputs],
            [t.dtype.as_datatype_enum for t in inputs],
    )
    # Save
    builder.add_meta_graph_and_variables(sess, ["bert-uncased"])
builder.save()

print("Exported SavedModel:", out_dir)
