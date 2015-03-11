Huabot Brain
============

An Image Recognition System build on top of caffe deep learn framework.

Quick start
-----------

Please make sure install the requirements.

### install

    go get -v github.com/Lupino/huabot-brain
    go get -v github.com/Lupino/huabot-brain/tools/train_worker

### start server

    cd $GOPATH/src/github.com/Lupino/huabot-brain
    make # precompile and package javascript
    huabot-brain --gearmand=127.0.0.1:4730 --dbpath=dataset.db
    train_worker --gearmand=127.0.0.1:4730

    env GEARMAND_PORT=tcp://127.0.0.1:4730 python tools/predict_worker/main.py resources

### load datasets

    cd $GOPATH/src/github.com/Lupino/huabot-brain/tools/datasets
    python get_datasets.py

### open dashboard

Go to <http://127.0.0.1:3000>

Just click Train button to train the network.

### learn more

see [API.md](https://github.com/Lupino/huabot-brain/blob/master/API.md)

Requirements
------------

* [caffe](http://caffe.berkeleyvision.org/)
* [Gearman](http://gearman.org)
* [Python](http://python.org)
* [golang](http://golang.org)
