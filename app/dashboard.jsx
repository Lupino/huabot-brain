var Dashboard = React.createClass({
  contextTypes: {
    router: React.PropTypes.func.isRequired
  },

  loadTags () {
    var self = this;
    var query = this.context.router.getCurrentQuery();
    var max = query.max || '';
    var limit = Number(query.limit) || 20;
    jQuery.get('/api/tags/?max=' + max + '&limit=' + limit, data => {
      data.has_more = data.tags.length >= limit;
      data.lastTag = data.tags[data.tags.length - 1];
      self.setState(data);
    });

  },

  loadStatus () {
    var self = this;
    jQuery.get('/api/solve', data => {
      self.setState(data);
    });
  },

  handleSolver () {
    var self = this;
    jQuery.post('/api/solve', data =>
      self.setState({status: 'Solving', loss: 0, acc: 0})
    );
  },

  handleStopSolver () {
    var self = this;
    if (confirm("Are you sure stop the Solving?")) {
      jQuery.ajax({url: '/api/solve', method: 'DELETE'}).done(() =>
        self.setState({status: 'Solved'})
      );
    }
  },

  handleClickTag (evt) {
    var elem = evt.target;
    var action = elem.getAttribute("action");
    var tagId = Number(elem.getAttribute("data-id"));
    var self = this;

    if (action === 'edit') {
      var tagName = prompt('Enter tag:');
      jQuery.post('/api/tags/' + tagId, {tag: tagName}, data => {
        if (data.err) {
          alert(data.err);
          return;
        }

        self.cache.tags = self.cache.tags.map(tag => {
          if (tag.tag_id === tagId) {
            tag.name = tagName;
          }
          return tag;
        });

        self.setState({updateTag: tagId, tags: []});
      });
    } else if (action === 'delete') {
      if (confirm("Are you sure?")) {
        jQuery.ajax({url: '/api/tags/' + tagId, method: 'DELETE'}).done(() => {
          self.cache.tags = self.cache.tags.filter(tag => tag.tag_id !== tagId);
          self.setState({removeTag: tagId, tags: []});
        }).fail(() =>
          alert("Error: Not Found.")
        );
      }
    }
  },

  getInitialState () {
    this.cache = this.cache || {};
    return {tags: [], status: 'Solved', acc: 0, loss: 0,
            removeTag: false, updateTag: false,
            has_more: false, lastTag: null};
  },

  shouldLoadTags () {
    var path = this.context.router.getCurrentPath();
    if (this.cache.path !== path) {
      this.cache.path = path;
      return true;
    }
    return false;
  },

  shouldCleanTags () {
    if (!this.context.router.getCurrentQuery().max) {
      return true;
    }
    return false;
  },

  cleanTags () {
    this.cache.tags = [];
  },

  componentDidMount () {
    this.cache.tags = this.cache.tags || [];
    this.loadStatus();
    this.componentDidUpdate();
  },

  componentDidUpdate () {
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

  renderSolver () {
    var btn = <Button bsStyle="primary" bsSize="xsmall" onClick={this.handleSolver}> Solve </Button>;
    if (this.state.status == "Solving") {
      btn = <Button bsStyle="danger" bsSize="xsmall" onClick={this.handleStopSolver}> Stop </Button>
    }
    return (
      <div className="solver">
        <Panel header="Solver status" bsStyle="info">
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
      </div>
    );
  },

  render () {
    var tags = this.state.tags || [];
    var loadMore;
    var {router} = this.context;
    if (this.state.has_more && this.state.lastTag) {
      var query = router.getCurrentQuery();
      query = query || {};
      query.max = this.state.lastTag.tag_id;
      loadMore = (
        <div className="load-more">
          <ButtonLink bsStyle="info" bsSize="large"
              params={router.getCurrentParams()} to="dashboard" query={query} block>Load More...</ButtonLink>
        </div>
      );
    }
    if (this.cache.tags && this.cache.tags.length > 0) {
      var oldLastTag = this.cache.tags[this.cache.tags.length - 1];
      var lastTag = this.state.lastTag;
      if (lastTag && oldLastTag.tag_id !== lastTag.tag_id) {
        this.cache.tags = this.cache.tags.concat(tags);
      }
    } else {
      this.cache.tags = tags;
    }
    var elems = this.cache.tags.map(tag => {
      return (
        <tr key={tag.tag_id}>
          <td>{tag.tag_id}</td>
          <td><Link to="datasets" params={{dataType: 'all'}} query={{tag: tag.name}}>{tag.name}</Link></td>
          <td>{tag.train_count}</td>
          <td>{tag.test_count}</td>
          <td>
            <Button bsSize="xsmall" data-id={tag.tag_id} action="edit">Edit</Button>
            &nbsp;
            <Button bsStyle="danger" bsSize="xsmall" data-id={tag.tag_id} action="delete">Delete</Button>
          </td>
        </tr>
      );
    });
    return (
      <div className="dashboard">
        {this.renderSolver()}
        <h2 className="sub-header">Tags</h2>
        <Table striped bordered condensed hover onClick={this.handleClickTag}>
          <thead>
            <tr>
              <th>#</th>
              <th>Name</th>
              <th>Train</th>
              <th>Test</th>
              <th width={100}></th>
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
