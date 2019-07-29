from __future__ import print_function

import argparse
import os
import os.path
import urllib.request  # TODO py2/3 compat
import zipfile

parser = argparse.ArgumentParser()
parser.add_argument("model", help="Name of the model to download")
parser.add_argument("output_path",
                    help="Path to save model to, model name will be appended")

parser.add_argument("-f", "--force", action='store_true',
                    help="Force a download, even if model exists")
parser.add_argument("-k", "--keep", action='store_true',
                    help="Keep the downloaded archive")

BERT_PRETRAINED_MODEL_URLS = {
    'bert-base-uncased': "https://storage.googleapis.com/bert_models/2018_10_18/uncased_L-12_H-768_A-12.zip",  # noqa
    'bert-large-uncased': "https://storage.googleapis.com/bert_models/2018_10_18/uncased_L-24_H-1024_A-16.zip",  # noqa
    'bert-base-cased': "https://storage.googleapis.com/bert_models/2018_10_18/cased_L-12_H-768_A-12.zip",  # noqa
    'bert-large-cased': "https://storage.googleapis.com/bert_models/2018_10_18/cased_L-24_H-1024_A-16.zip",  # noqa
    'bert-base-multilingual-cased': "https://storage.googleapis.com/bert_models/2018_11_23/multi_cased_L-12_H-768_A-12.zip",  # noqa
    'bert-base-chinese': "https://storage.googleapis.com/bert_models/2018_11_03/chinese_L-12_H-768_A-12.zip",  # noqa
    'bert-large-uncased-wwm': "https://storage.googleapis.com/bert_models/2019_05_30/wwm_uncased_L-24_H-1024_A-16.zip",  # noqa
    'bert-large-cased-wwm': "https://storage.googleapis.com/bert_models/2019_05_30/wwm_cased_L-24_H-1024_A-16.zip",  # noqa
}


def download(model, output_path, force=False, keep=False):
    url = BERT_PRETRAINED_MODEL_URLS.get(model, None)
    if not url:
        print("Invalid Model Name:", model)
        for k in BERT_PRETRAINED_MODEL_URLS.keys():
            print("Valid Model Names:")
            print("\t", k)
            exit(1)
    path = os.path.join(output_path, model)
    if os.path.exists(path):
        if not force:
            print("Model Already Exists:", path)
            print("If desired, use --force to overwrite")
            return
    os.makedirs(path, exist_ok=True)  # TODO (py3.2+) py2 compat
    print("Downloading and extracting model:", model)
    print("Downloading archive:", url)
    archive = os.path.join(path, os.path.basename(url))
    urllib.request.urlretrieve(url, archive)
    print("Extracting archive to path:", path)
    _unzip(archive, path)
    if not keep:
        print("Removing archive:", archive)
        os.remove(archive)
    print("Extracted Model", model, "to", path)
    print("Done.")


def _unzip(archive, dst):
    """ unzips archive and flattens directory structure """
    with zipfile.ZipFile(archive) as zf:
        for zi in zf.infolist():
            if zi.filename[-1] == '/':  # skip dir
                continue
            zi.filename = os.path.basename(zi.filename)
            zf.extract(zi, dst)


if __name__ == '__main__':
    args = parser.parse_args()
    model = args.model
    output_path = args.output_path
    force = args.force
    keep = args.keep
    download(model, output_path, force, keep)
