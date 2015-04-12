#!/usr/bin/env python

import os
import time
import numpy as np
import caffe
import cStringIO as StringIO
import json
import urllib2
from urlparse import urlparse
from caffe.io import resize_image
import logging
from bottle import request, response, Bottle, run
# logging.basicConfig(level=logging.DEBUG)

caffe.set_mode_cpu()
MAX_PREDICT_LENGTH = 5

RAW_SCALE = 255.

def load_binaryproto(fn):
    blob = caffe.proto.caffe_pb2.BlobProto()
    data = open(fn, 'rb').read()
    blob.ParseFromString(data)
    arr = np.array( caffe.io.blobproto_to_array(blob) )
    return arr[0]


class Classifier(object):
    def __init__(self, resoursesPath):
        mean_file = resoursesPath + "/mean.binaryproto"
        model_def_file = resoursesPath + "/deploy.prototxt"
        pretrained_model_file = resoursesPath + "/models/huabot-brain.caffemodel"

        mean=load_binaryproto(mean_file)

        self.net = caffe.Classifier(
            model_def_file, pretrained_model_file,
            image_dims=(256, 256),
            raw_scale=RAW_SCALE,
            channel_swap=(2, 1, 0)
        )

        in_shape = self.net.transformer.inputs[self.net.inputs[0]]
        if mean.shape[1:] != in_shape[2:]:
            mean = caffe.io.resize_image(mean.transpose((1,2,0)), in_shape[2:]).transpose((2,0,1))

        self.net.transformer.set_mean(self.net.inputs[0], mean)

    def classify_image(self, image):
        try:
            starttime = time.time()
            scores = self.net.predict([image], oversample=True).flatten()
            endtime = time.time()

            indices = (-scores).argsort()[:MAX_PREDICT_LENGTH]
            meta = [{'id':i, 'score': float(scores[i])} for i in indices]
            return (True, meta, endtime - starttime)

        except Exception as err:
            logging.exception(err)
            return (False, 'Something went wrong when classifying the '
                          'image. Maybe try another one?')

clf = Classifier('resourses')

app = Bottle()

@app.post("/api/predict/url")
def predict_url():
    response.set_header('content-type', 'application/json')
    img_url = request.forms.img_url.strip()
    if not img_url:
        return json.dumps({"err": "img_url is required."})
    rsp = urllib2.urlopen(url, timeout=10)
    data = rsp.read()
    image = caffe.io.load_image(StringIO.StringIO(data))
    result = clf.classify_image(image)
    if result[0]:
        result = {'bet_result': result[1], 'time': result[2]}
    else:
        result = {'err': result[1]}

    return json.dumps(result)

def main(script, arg1, val1):
    if arg1 == '--host':
        host = urlparse(val1).netloc
    else:
        host = "127.0.0.1:3001"

    host = host.split(":")
    host = host[0]
    port = "80"
    if len(host) == 2:
        port = host[1]

    run(app, host=host, port=port)

if __name__ == "__main__":
    import sys
    main(*sys.argv)
