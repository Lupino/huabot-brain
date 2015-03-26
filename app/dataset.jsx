var Dataset = React.createClass({
  handleDelete () {
    var self = this;
    jQuery.ajax({
      method: 'DELETE',
      url: '/api/datasets/' + this.props.data.dataset_id
    }).done(() => {
      self.props.onRemove(self.props.data.dataset_id);
      self.props.onRequestHide();
    });
  },
  render () {
    var dataset = this.props.data;
    var ext = FILE_EXTS[dataset.file.type] || '.jpg';
    return (
      <Modal {...this.props} title={dataset.tag.name} animation={false}>
        <div className="modal-body">
          <div className="img" data-id={dataset.dataset_id}>
            <div className="dataset">
              <img src={"/upload/" + dataset.file.key + ext} />
            </div>
          </div>
        </div>
        <div className="modal-footer">
          <Button bsStyle="danger" onClick={this.handleDelete}>Delete</Button>
          <Button onClick={this.props.onRequestHide}>Close</Button>
        </div>
      </Modal>
    )
  }
});


var Datasets = React.createClass({
  contextTypes: {
    router: React.PropTypes.func.isRequired
  },

  waterfall () {
    if (this.cache.datasets.length === 0) {
      return;
    }
    jQuery("#waterfall").waterfall({
      selector: ".dataset",
    });
  },

  loadDatasets () {
    var self = this;
    var {router} = this.context;
    var query = router.getCurrentQuery();
    var params = router.getCurrentParams();
    var max = query.max || '';
    var limit = Number(query.limit) || 50;
    var tag = query.tag || '';
    var dataType = params.dataType || 'all';
    if (this.cache.dataType && this.cache.dataType !== dataType) {
      max = '';
    }
    this.cache.dataType = dataType;
    jQuery.get('/api/datasets/?max=' + max + '&limit=' + limit + '&data_type=' + dataType + '&tag=' + tag, function(data) {
      data.has_more = data.datasets.length >= limit;
      data.lastDataset = data.datasets[data.datasets.length - 1];
      self.setState(data);
    });
  },

  handleRemoveDataset (datasetId) {
    this.cache.datasets = this.cache.datasets.filter(d => d.dataset_id !== datasetId);
    this.setState({removeDataset: datasetId, datasets: []});
  },

  getInitialState () {
    this.cache = this.cache || {};
    return {datasets: [], lastDataset: null, has_more: false};
  },

  shouldLoadDatasets () {
    var path = this.context.router.getCurrentPath();
    if (this.cache.path !== path) {
      this.cache.path = path;
      return true;
    }
    return false;
  },

  shouldCleanDatasets () {
    var {router} = this.context;
    var query = router.getCurrentQuery();
    var pathname = router.getCurrentPathname();
    if (this.cache.pathname !== pathname || this.cache.tag !== query.tag || !query.max) {
      this.cache.pathname = pathname;
      this.cache.tag = query.tag;
      this.cache.scroll = true;
      return true;
    }
    return false;
  },

  cleanDatasets () {
    this.cache.datasets = [];
  },

  componentDidMount () {
    this.cache.datasets = this.cache.datasets || [];
    this.componentDidUpdate();
  },

  componentDidUpdate () {
    if (this.shouldLoadDatasets()) {
      if (this.shouldCleanDatasets()) {
        this.cleanDatasets();
      }
      this.loadDatasets();
    } else {
      if (this.cache.scroll) {
        this.cache.scroll = false;
        window.scroll(0, 0);
      }
      this.waterfall();
    }
  },

  render () {
    var self = this;
    var {router} = this.context;
    var datasets = this.state.datasets || [];
    var loadMore;
    if (this.state.has_more && this.state.lastDataset) {
      var query = router.getCurrentQuery();
      query = query || {};
      query.max = this.state.lastDataset.dataset_id;
      loadMore = (
        <div className="load-more">
          <ButtonLink bsStyle="info" bsSize="large"
              params={router.getCurrentParams()} to="datasets" query={query} block>Load More...</ButtonLink>
        </div>
      );
    }
    if (this.cache.datasets && this.cache.datasets.length > 0) {
      var oldLastDataset = this.cache.datasets[this.cache.datasets.length - 1];
      var lastDataset = this.state.lastDataset;
      if ( lastDataset && oldLastDataset.dataset_id !== lastDataset.dataset_id) {
        this.cache.datasets = this.cache.datasets.concat(datasets);
      }
    } else {
      this.cache.datasets = datasets;
    }
    var elems = this.cache.datasets.map(dataset => {
      var width = 192;
      var height = width / dataset.file.width * dataset.file.height;
      if (height > 600) {
        height = 600;
      }
      var ext = FILE_EXTS[dataset.file.type] || '.jpg';
      return (
        <ModalTrigger modal={<Dataset data={dataset} title={dataset.tag.name}
            onRemove={self.handleRemoveDataset} />} key={dataset.dataset_id}>
          <div className="dataset" data-id={dataset.dataset_id}>
            <div className="file" style={{width: width, height: height}}>
              <img src={"/upload/" + dataset.file.key + ext} />
            </div>
            <div className="tag">{dataset.tag.name}</div>
          </div>
        </ModalTrigger>
      );
    });
    return (
      <div className="datasets">
        <div id="waterfall">
          {elems}
        </div>
        {loadMore}
      </div>
    );
  }
});

var NewDataset = React.createClass({

  handleClick () {
    $(this.refs.file.getDOMNode()).click();
  },

  getInitialState () {
    return {
      file: null,
      dataset: null
    };
  },

  componentDidMount  () {
    var self = this;
    $(".fileForm").ajaxForm({
      success(result) {self.setState(result);},
      error() { alert("Please choose an image...");}
    });
  },

  handleFile () {
    $(this.refs.fileForm.getDOMNode()).submit();
  },

  handleToggle () {
    this.props.onRequestHide();
    if (this.state.dataset) {
      window.location.reload();
    }
  },

  handleDataTypeClick (evt) {
    var dataType = evt.target.value;
    this.setState({data_type: dataType});
  },

  handleSave () {
    var self = this;
    var file_id = this.state.file.file_id;
    var tag = this.refs.tag.getValue();
    var desc = this.refs.desc.getValue();
    var dataType = this.state.data_type;
    if (!tag) {
      alert("Tag is required.");
    }
    jQuery.post("/api/datasets",
        {tag: tag, file_id: file_id, description: desc, data_type: dataType},
        data => self.setState(data));
  },

  render () {
    var action = this.props.action || '/api/upload';
    var fileForm, saveBtn, mainBody, dataType;

    if (this.state.dataset) {
      var ext = FILE_EXTS[this.state.dataset.file.type] || '.jpg';
      var mainBody = (
        <div className="img" data-id={this.state.dataset.dataset_id}>
          <div className="dataset">
            <img src={"/upload/" + this.state.dataset.file.key + ext} />
          </div>
        </div>
      );
    } else if (this.state.file) {
      saveBtn = <Button bsStyle="primary" onClick={this.handleSave}>Save</Button>;
      var height, width;
      var boxStyle = {};
      if (this.state.file.width > this.state.file.height) {
        width = 136;
        height = this.state.file.height / this.state.file.width * width;
      } else {
        height = 136;
        width = this.state.file.width / this.state.file.height * height;
      }
      boxStyle.paddingTop = (136 - height) / 2;
      boxStyle.paddingLeft = (136 - width) / 2;

      var ext = FILE_EXTS[this.state.file.type] || '.jpg';
      mainBody = (
        <Row className="new-dataset">
          <Col xs={6} md={4}>
            <Panel>
              <div className="imgBox" style={boxStyle}>
                <img src={"/upload/" + this.state.file.key + ext} width={width} height={height} />
              </div>
             </Panel>
          </Col>
          <Col xs={12} md={8}>
            <form ref="datasetForm" encType="multipart/form-data">
              <Input ref="tag" type="text" label="Tag:" />
              <Input ref="desc" type="textarea" label="Description:" className="desc" />
            </form>
          </Col>
        </Row>
      );

      dataType = (
        <div className="dataType">
          <Row>
            <Col xs={6} md={4}>
              <Input type="radio" label="Candidate" name="data_type" value={0} onClick={this.handleDataTypeClick}  />
            </Col>
            <Col xs={6} md={4}>
              <Input type="radio" label="Train" name="data_type" value={1} onClick={this.handleDataTypeClick}  />
            </Col>
            <Col xs={6} md={4}>
              <Input type="radio" label="Val"  name="data_type" value={2} onClick={this.handleDataTypeClick} />
            </Col>
          </Row>
        </div>

      );

    } else {
      fileForm = (
        <div className="fileForm">
          <Button bsStyle="primary" bsSize="large" block onClick={this.handleClick}>Choose file...</Button>
          <form ref="fileForm" encType="multipart/form-data" method="POST" action={action}>
            <input ref="file" type="file" name="file" onChange={this.handleFile} />
          </form>
        </div>
      );
    }
    return (
        <Modal bsStyle="primary" title="Add new Dataset" onRequestHide={this.handleToggle}>
          <div className="modal-body">
            {fileForm}
            {mainBody}
          </div>
          <div className="modal-footer">
            {dataType}
            <Button onClick={this.handleToggle}>{this.state.dataset ? "Close" : "Cancel"}</Button>
            {saveBtn}
          </div>
        </Modal>
      );
  }
});
