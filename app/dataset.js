var Dataset = React.createClass({
  render: function() {
    var dataset = this.props.data;
    return (
      <Modal {...this.props} title="View" animation={false}>
        <div className="modal-body">
          <div className="img" data-id={dataset.dataset_id}>
            <div className="dataset">
              <img src={"/upload/" + dataset.file.key} />
              <div className="tag">{dataset.tag.name}</div>
            </div>
          </div>
        </div>
        <div className="modal-footer">
          <Button onClick={this.props.onRequestHide}>Close</Button>
        </div>
      </Modal>
    )
  }
});


var Datasets = React.createClass({
  mixins: [State],

  waterfall: function() {
    if (this.cache.datasets.length === 0) {
      return;
    }
    jQuery("#waterfall").waterfall({
      selector: ".dataset",
    });
  },

  loadDatasets: function() {
    var self = this;
    var query = this.getQuery();
    var params = this.getParams();
    var max = query.max || '';
    var limit = Number(query.limit) || 50;
    var tag = query.tag || '';
    this.limit = limit;
    var dataType = params.dataType || 'all';
    jQuery.get('/api/datasets/?max=' + max + '&limit=' + limit + '&data_type=' + dataType + '&tag=' + tag, function(data) {
      self.setState(data);
    });
  },

  getInitialState: function() {
    this.cache = this.cache || {};
    return {datasets: []};
  },

  shouldLoadDatasets: function() {
    var path = this.getPath();
    if (this.cache.path !== path) {
      this.cache.path = path;
      return true;
    }
    return false;
  },

  shouldCleanDatasets: function() {
    var pathname = this.getPathname();
    var query = this.getQuery();
    if (this.cache.pathname !== pathname || this.cache.tag !== query.tag || !query.max) {
      this.cache.pathname = pathname;
      this.cache.tag = query.tag;
      this.cache.scroll = true;
      return true;
    }
    return false;
  },

  cleanDatasets: function() {
    this.cache.datasets = [];
  },

  componentDidMount: function() {
    this.cache.datasets = this.cache.datasets || [];
    this.componentDidUpdate();
  },

  componentDidUpdate: function() {
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

  render: function() {
    var datasets = this.state.datasets || [];
    var loadMore;
    if (datasets.length >= this.limit) {
      var query = this.getQuery();
      query = query || {};
      query.max = datasets[datasets.length - 1].dataset_id;
      loadMore = (
        <div className="load-more">
          <ButtonLink bsStyle="info" bsSize="large"
              params={this.getParams()} to="datasets" query={query} block>Load More...</ButtonLink>
        </div>
      );
    }
    if (this.cache.datasets && this.cache.datasets.length > 0) {
      var oldLastDataset = this.cache.datasets[this.cache.datasets.length - 1];
      var lastDataset = datasets[datasets.length - 1];
      if (oldLastDataset.dataset_id !== lastDataset.dataset_id) {
        this.cache.datasets = this.cache.datasets.concat(datasets);
      }
    } else {
      this.cache.datasets = datasets;
    }
    var elems = this.cache.datasets.map(function(dataset) {
      var width = 192;
      var height = width / dataset.file.width * dataset.file.height;
      if (height > 600) {
        height = 600;
      }
      return (
        <ModalTrigger modal={<Dataset data={dataset} title={dataset.tag.name} />}>
          <div className="dataset" data-id={dataset.dataset_id}>
            <div className="file" style={{width: width, height: height}}>
              <img src={"/upload/" + dataset.file.key} />
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
