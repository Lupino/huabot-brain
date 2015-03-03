POST /api/datasets
------------------

Add a Dataset.

|param|type|required|desc|
|:--|:--:|:--:|:--|
|tag|str|true||
|data_type|int|true|TRAIN: 1, VAL: 2|
|file|blob|true|image data|


POST /api/datasets/:dataset_id
------------------------------

Update dataset datatype.

|param|type|required|desc|
|:--|:--:|:--:|:--|
|data_type|int|true|TRAIN: 1, VAL: 2|


GET /api/datasets
-----------------
Get the datasets list.

|param|type|required|desc|
|:--|:--:|:--:|:--|
|tag|str|false||
|data_type|str|false|train, val, default: all|
|max|int|false|max dataset_id|
|limit|int|false|max 100, default 10|

GET /api/datasets/:dataset_id
-----------------------------
Get a dataset.

DELETE /api/datasets/:dataset_id
--------------------------------

Delete a dataset.


POST /api/tags
------------------

Add a Dataset.

|param|type|required|desc|
|:--|:--:|:--:|:--|
|tag|str|true|the dataset label.|


POST /api/tags/:tag_id
------------------------------

Update tag datatype.

|param|type|required|desc|
|:--|:--:|:--:|:--|
|tag|str|true|the dataset label.|


GET /api/tags
-----------------
Get the tags list.

|param|type|required|desc|
|:--|:--:|:--:|:--|
|max|int|false|max tag_id|
|limit|int|false|max 100, default 10|

GET /api/tags/:tag_id
-----------------------------
Get a tag.

DELETE /api/tags/:tag_id
--------------------------------

Delete a tag.

GET /api/tags/hint
------------------

Get a word hint.

|param|type|required|desc|
|:--|:--:|:--:|:--|
|word|str|true|the word which want hint|

POST /api/train
---------------

Train the network.

GET /api/train
--------------

Get the train status

DELETE /api/train
-----------------

Stop train

GET /api/train.txt
------------------

Build file train.txt

GET /api/val.txt
----------------

Build file val.txt

GET /api/loss.png
-----------------

Draw train loss

GET /api/acc.png
----------------
Draw test acc

POST /api/predict
-----------------

Predict a image tags.

|param|type|required|desc|
|:--|:--:|:--:|:--|
|img_url|str|true|the image url.|

GET /api/proxy
--------------

Proxy to load image.

|param|type|required|desc|
|:--|:--:|:--:|:--|
|url|str|true|the image url.|

POST /api/upload

Upload a file.

|param|type|required|desc|
|:--|:--:|:--:|:--|
|file|blob|true|the image blob data.|
