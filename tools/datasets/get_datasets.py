#!/usr/bin/env python

import requests
import os
from io import BytesIO

TRAIN_LIMIT = 100
VAL_LIMIT = 50
TRAIN = 1
VAL = 2

BRAIN_ROOT = os.environ.get('BRAIN_ROOT', 'http://127.0.0.1:3000/api')

def get_pins(uri):
    api = 'http://api.huaban.com/{}'.format(uri)
    url = api

    while True:
        rsp = requests.get(url, headers={
            'Accept': 'application/json',
            'User-Agent': 'Huabot/1.0'})
        pins = rsp.json().get('pins', [])
        if len(pins) == 0:
            break

        min_pin_id = pins[0]['pin_id']
        for pin in pins:
            if min_pin_id > pin['pin_id']:
                min_pin_id = pin['pin_id']

            yield pin

        url = api + '?limit=100&max=%s'%min_pin_id


def submit_pin(pin, tag, data_type=1):
    file_url = 'http://img.hb.aicdn.com/' + pin['file']['key'] + '_fw320'
    print('Submit dataset: %s (%s)'%(tag, file_url))
    try:
        res = requests.get(file_url, timeout=30)
        data = res.content
        f = BytesIO()
        f.write(data)
        f.seek(0, 0)
        res = requests.post(BRAIN_ROOT + '/datasets/',
                            files={"file": f},
                            data={'tag': tag, "data_type": data_type})
        return True
    except Exception as e:
        print(e)
        return False


def load_datasets(uri, tag):
    print('Load datasets: ' + tag)
    pins = get_pins(uri)
    count = 0
    for pin in pins:
        if submit_pin(pin, tag, VAL):
            count += 1
        if count > VAL_LIMIT:
            break

    count = 0
    for pin in pins:
        if submit_pin(pin, tag, TRAIN):
            count += 1
        if count > TRAIN_LIMIT:
            break

def main():

    datasets = [dataset.split('\t')
                for dataset in open('datasets.txt', 'r').read().split('\n') if dataset]

    for dataset in datasets[:]:
        uri = dataset[1].strip()
        tag = dataset[0].strip()
        load_datasets(uri, tag)

if __name__ == "__main__":
    main()
