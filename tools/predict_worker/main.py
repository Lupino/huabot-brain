import os
import time
import numpy as np
import caffe
import gear
import cStringIO as StringIO
import json
import urllib2
from caffe.io import resize_image
import logging
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

class Brain(object):
    def __init__(self, resoursesPath):
        self._clf = Classifier(resoursesPath)
        self._funcs = {}
        self._worker = gear.Worker('huaban-brain')

    def _add_func(self, func_name, callback):
        self._worker.registerFunction(func_name)
        self._funcs[func_name] = callback

    def add_server(self, host='localhost', port=4730):
        self._worker.addServer(host, port)

    def process(self):
        self._add_func('CAFFE:PREDICT', self.classify_image)
        self._add_func('CAFFE:PREDICT:URL', self.classify_image_url)
        while 1:
            job = self._worker.getJob()
            func = self._funcs.get(job.name)
            if func:
                try:
                    func(job)
                except Exception as e:
                    job.sendWorkComplete(json.dumps({'err': str(e)}))
                    print('process %s error: %s'%(job.name, e))

    def classify_image(self, job):
        self._classify_image(job, job.arguments)

    def classify_image_url(self, job):
        url = job.arguments

        rsp = urllib2.urlopen(url, timeout=10)
        data = rsp.read()
        self._classify_image(job, data)

    def _classify_image(self, job, data):
        image = caffe.io.load_image(StringIO.StringIO(data))
        result = self._clf.classify_image(image)
        if result[0]:
            result = {'bet_result': result[1], 'time': result[2]}
        else:
            result = {'err': result[1]}
        print(result)
        job.sendWorkComplete(json.dumps(result))


def main(scripts, resoursesPath='resourses'):
    GEARMAND_PORT = os.environ.get('GEARMAND_PORT',
                               'tcp://127.0.0.1:4730')[6:].split(':')

    brain = Brain(resoursesPath)
    brain.add_server(GEARMAND_PORT[0], int(GEARMAND_PORT[1]))

    print("brain process")
    brain.process()

if __name__ == "__main__":
    import sys
    main(*sys.argv)
