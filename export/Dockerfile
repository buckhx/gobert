FROM tensorflow/tensorflow:1.13.2-py3

WORKDIR /workspace/export

COPY . .

ENTRYPOINT ["./entrypoint.sh"]
CMD ["export_embedding.py", "--download=bert-base-uncased", "/var/bert/model", "/var/bert/export/embedding"]
