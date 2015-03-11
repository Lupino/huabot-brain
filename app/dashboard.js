var Dashboard = React.createClass({
  mixins: [State],

  loadTags: function() {
    var self = this;
    var query = this.getQuery();
    var max = query.max || '';
    var limit = Number(query.limit) || 20;
    this.limit = limit;
    jQuery.get('/api/tags/?max=' + max + '&limit=' + limit, function(data) {
      self.setState(data);
    });

  },

  loadStatus: function() {
    var self = this;
    jQuery.get('/api/train', function(data) {
      self.setState(data);
    });
  },

  handleTrain: function() {
    var self = this;
    jQuery.post('/api/train', function(data) {
      self.setState({status: 'training', loss: 0, acc: 0});
    });
  },

  handleStopTrain: function() {
    var self = this;
    if (confirm("Are you sure stop the training?")) {
      jQuery.ajax({url: '/api/train', method: 'DELETE'}, function(data) {
        self.setState({status: 'no train'});
      });
    }
  },

  getInitialState: function() {
    this.cache = this.cache || {};
    return {tags: [], status: 'no train', acc: 0, loss: 0};
  },

  shouldLoadTags: function() {
    var path = this.getPath();
    if (this.cache.path !== path) {
      this.cache.path = path;
      return true;
    }
    return false;
  },

  shouldCleanTags: function() {
    if (!this.getQuery().max) {
      return true;
    }
    return false;
  },

  cleanTags: function() {
    this.cache.tags = [];
  },

  componentDidMount: function() {
    this.cache.tags = this.cache.tags || [];
    this.loadStatus();
    this.componentDidUpdate();
  },

  componentDidUpdate: function() {
    if (this.shouldLoadTags()) {
      if (this.shouldCleanTags()) {
        this.cache.scroll = true;
        this.cleanTags();
      }
      this.loadTags();
    } else {
      if (this.cache.scroll) {
        this.cache.scroll = false;
        window.scroll(0, 0);
      }
    }
  },

  render: function() {
    var tags = this.state.tags || [];
    var loadMore;
    if (tags.length >= this.limit) {
      var query = this.getQuery();
      query = query || {};
      query.max = tags[tags.length - 1].tag_id;
      loadMore = (
        <div className="load-more">
          <ButtonLink bsStyle="info" bsSize="large"
              params={this.getParams()} to="dashboard" query={query} block>Load More...</ButtonLink>
        </div>
      );
    }
    if (this.cache.tags && this.cache.tags.length > 0) {
      var oldLastTag = this.cache.tags[this.cache.tags.length - 1];
      var lastTag = tags[tags.length - 1];
      if (oldLastTag.tag_id !== lastTag.tag_id) {
        this.cache.tags = this.cache.tags.concat(tags);
      }
    } else {
      this.cache.tags = tags;
    }
    var elems = this.cache.tags.map(function(tag) {
      return (
        <tr>
          <td>{tag.tag_id}</td>
          <td><Link to="datasets" params={{dataType: 'all'}} query={{tag: tag.name}}>{tag.name}</Link></td>
          <td>{tag.train_count}</td>
          <td>{tag.test_count}</td>
        </tr>
      );
    });
    var btn = <Button bsStyle="primary" bsSize="xsmall" onClick={this.handleTrain}> Train </Button>;
    if (this.state.status == "training") {
      btn = <Button bsStyle="danger" bsSize="xsmall" onClick={this.handleStopTrain}> Stop </Button>
    }
    return (
      <div className="dashboard">
        <Panel header="Train status" bsStyle="info">
          <Row>
            <Col xs={6}>
              <img src="/api/loss.png" />
            </Col>
            <Col xs={6}>
              <img src="/api/acc.png" />
            </Col>
          </Row>
        </Panel>
        <Panel>
          <Row>
            <Col xs={6} md={4}>Loss: {this.state.loss}</Col>
            <Col xs={6} md={4}>Accurancy: {this.state.acc}</Col>
            <Col xs={6} md={4}>
              <Row>
                <Col xs={12} md={8}> Status: {this.state.status} </Col>
                <Col xs={6} md={4}> {btn} </Col>
              </Row>
            </Col>
          </Row>
        </Panel>
        <h2 class="sub-header">Tags</h2>
        <Table striped bordered condensed hover>
          <thead>
            <tr>
              <th>#</th>
              <th>Name</th>
              <th>Train</th>
              <th>Test</th>
            </tr>
          </thead>
          <tbody>
            {elems}
          </tbody>
        </Table>
        {loadMore}
      </div>
    );
  }
});

