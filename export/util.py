import os.path
import shutil

import tensorflow as tf


def export(model_path, export_path, transfer_func,
           method_name="bert/tuned/predict", tags=["bert-tuned"],
           sig_name="classifier", vocab_file="vocab.txt"):
    inputs, outputs = transfer_func()
    inputs = {k: tf.saved_model.utils.build_tensor_info(t)
              for k, t in inputs.items()}
    outputs = {k: tf.saved_model.utils.build_tensor_info(t)
               for k, t in outputs.items()}
    saver = tf.train.Saver()
    sig = tf.saved_model.signature_def_utils.build_signature_def(
        inputs=inputs,
        outputs=outputs,
        method_name=method_name,
    )
    with tf.Session() as sess:
        # TODO investigate optimize_graph
        checkpoint = tf.train.latest_checkpoint(model_path)
        if checkpoint is None:  # pretrained doesn't have latest
            checkpoint = os.path.join(model_path, "bert_model.ckpt")
        saver.restore(sess, checkpoint)
        b = tf.saved_model.builder.SavedModelBuilder(export_path)
        b.add_meta_graph_and_variables(sess, tags, {sig_name: sig},
                                       clear_devices=True)
        b.save()
        if vocab_file:
            shutil.copyfile(os.path.join(model_path, vocab_file),
                            os.path.join(export_path, vocab_file))
