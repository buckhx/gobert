import tensorflow as tf
from tensorflow.python.tools.optimize_for_inference_lib import optimize_for_inference


def export(model_path, export_path, transfer_func,
           method_name="bert/tuned/predict", tags=["bert-tuned"],
           sig_name="classifier"):
    inputs, outputs = transfer_func()
    inputs = {k: tf.saved_model.utils.build_tensor_info(t) for k, t in inputs.items()}
    outputs = {k: tf.saved_model.utils.build_tensor_info(t) for k, t in outputs.items()}
    saver = tf.train.Saver()
    sig = tf.saved_model.signature_def_utils.build_signature_def(
        inputs=inputs,
        outputs=outputs,
        method_name=method_name,
    )
    with tf.Session() as restore:
        # TODO investigate optimize_graph
        checkpoint = tf.train.latest_checkpoint(model_path)
        if checkpoint is None: # untuned doesn't have latest
            checkpoint = model_path+"/bert_model.ckpt"
        saver.restore(restore, checkpoint)
    with tf.Session() as sess:
        b = tf.compat.v1.saved_model.builder.SavedModelBuilder(export_path)
        g = tf.get_default_graph().as_graph_def()
        sess.run(tf.global_variables_initializer())
        g = optimize_for_inference(g,
            [k for k in inputs.keys()],
            [k for k in outputs.keys()],
            [t.dtype for t in inputs.values()],
        )
        b.add_meta_graph_and_variables(sess, tags, {sig_name: sig},
                clear_devices=True)
        b.save()

