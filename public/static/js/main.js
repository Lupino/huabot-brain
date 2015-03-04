var Navbar = ReactBootstrap.Navbar;
var Nav = ReactBootstrap.Nav;
var NavItemLink = ReactRouterBootstrap.NavItemLink;
var ButtonLink = ReactRouterBootstrap.ButtonLink;
var ModalTrigger = ReactBootstrap.ModalTrigger;
var Modal = ReactBootstrap.Modal;
var Button = ReactBootstrap.Button;
var Input = ReactBootstrap.Input;
var Grid = ReactBootstrap.Grid;
var Row = ReactBootstrap.Row;
var Col = ReactBootstrap.Col;
var Well = ReactBootstrap.Well;
var Panel = ReactBootstrap.Panel;
var ListGroup = ReactBootstrap.ListGroup;
var ListGroupItem = ReactBootstrap.ListGroupItem;
var Table = ReactBootstrap.Table;
var State = ReactRouter.State;
var Navigation = ReactRouter.Navigation;
var RouteHandler = ReactRouter.RouteHandler;
var Route = ReactRouter.Route;
var Router = ReactRouter.Router;
var Link = ReactRouter.Link;
var DefaultRoute = ReactRouter.DefaultRoute;
var NotFoundRoute = ReactRouter.NotFoundRoute;

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

var SearchForm = React.createClass({
  mixins: [State],

  getInitialState: function() {
    var query = this.getQuery();
    this.cache = this.cache || {};

    return {
      value: query.tag || '',
      tags: []
    };
  },

  getHint: function(word) {
    var self = this;
    jQuery.get('/api/tags/hint?word=' + word, function(data) {
      self.setState(data);
    });
  },

  shouldCleanHint: function() {
    var path = this.getPath();
    if (this.cache.path === path) {
      return false;
    }
    this.cache.path = path;
    return true;

  },

  cleanHint: function() {
    this.setState({tags: []});
  },

  componentDidUpdate: function() {
    if (this.cache.changed) {
      if (this.shouldCleanHint()) {
        this.cache.changed = false;
        this.cleanHint();
      }
      return;
    }
    var query = this.getQuery();
    if (query.tag && query.tag !== this.state.value) {
      this.setState({
        value: query.tag
      });
    } else if (!query.tag && this.state.value) {
      this.setState({value: ''});
    }
    this.cache.path = this.getPath();
  },

  handleChange: function() {
    this.cache.changed = true;
    this.getHint(this.refs.tag.getValue());
    this.setState({
      value: this.refs.tag.getValue()
    });
  },

  handleSubmit: function(evt) {
    evt.preventDefault();
    this.cache.changed = false;
    this.props.onSubmit(this.state.value);
  },

  handleListClick: function(eventKey, href, target) {
    this.cache.changed = false;
    this.props.onSubmit(target);
    this.setState({tags: [], value: target});
  },

  render: function() {
    var list = this.state.tags.map(function(tag) {
      return <ListGroupItem target={tag.name}> {tag.name}</ListGroupItem>;
    });
    return (
      <form className="navbar-form navbar-right" onSubmit={this.handleSubmit}>
        <Input type="text"
          name="tag"
          value={this.state.value}
          ref="tag"
          placeholder="Search..."
          onChange={this.handleChange} />

        <div className="hint">
          <ListGroup onClick={this.handleListClick}>
            {list}
          </ListGroup>
        </div>
      </form>
    );
  }
});

var App = React.createClass({
  mixins: [State, Navigation],

  getInitialState: function() {
    var tag = this.getQuery().tag || '';
    return {tag: tag}
  },

  handleSubmit: function(tag) {
    var query = this.getQuery();
    var params = this.getParams();
    params.dataType = params.dataType || 'all';
    query.tag = tag;
    delete query.max;
    var href = this.makeHref('datasets', params, query);
    window.location.href = href;
    this.setState({tag: tag || ''});
  },

  render: function() {
    return (
      <div className="app-main">
        <Navbar fixedTop inverse fluid brand="Huabot Brain">
          <Nav right>
            <NavItemLink to="dashboard">Dashboard</NavItemLink>
          </Nav>
          <SearchForm onSubmit={this.handleSubmit} />
        </Navbar>
        <Grid fluid>
          <Row>
            <Col sm={3} md={2} className="sidebar">
              <Nav>
                <NavItemLink to="datasets" params={{dataType: 'all'}}>All Data</NavItemLink>
                <NavItemLink to="datasets" params={{dataType: 'train'}}>Train Data</NavItemLink>
                <NavItemLink to="datasets" params={{dataType: 'val'}}>Val Data</NavItemLink>
                <NavItemLink to="demo">DEMO</NavItemLink>
              </Nav>
            </Col>
            <Col sm={9} smOffset={3} md={10} mdOffset={2}>
              <RouteHandler />
            </Col>
          </Row>
        </Grid>
      </div>
    );
  }
});

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
    })
  },

  handleTrain: function() {
    var self = this;
    jQuery.post('/api/train', function(data) {
      self.setState({status: 'training', loss: 0, acc: 0})
    });
  },

  handleStopTrain: function() {
    var self = this;
    if (confirm("Are you sure stop the training?")) {
      jQuery.ajax({url: '/api/train', method: 'DELETE'}, function(data) {
        self.setState({status: 'no train'})
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

var DEMO = React.createClass({
  mixins: [State, Navigation],

  predict: function(img_url) {
    var self = this;
    jQuery.post('/api/predict', {img_url: img_url}, function(data) {
      self.setState(data);
    });

  },

  getInitialState: function() {
    this.cache = {};
    var query = this.getQuery();
    var img_url = query.img_url || '';
    return {bet_result: [], time: null, err: null};
  },

  componentDidMount: function() {
    this.cache.imgUrl = 'http://img.hb.aicdn.com/d3a5039f151ddf451b95ea1b9b7e6af73f189c3a5f23-MMFBG3_fw320';
    this.predict(this.cache.imgUrl);
    this.componentDidUpdate();
  },

  componentDidUpdate: function() {
    var query = this.getQuery();
    if (query.img_url && query.img_url !== this.cache.imgUrl) {
      this.cache.imgUrl = query.img_url;
      this.predict(this.cache.imgUrl);
    }
  },

  handleSubmit: function(evt) {
    evt.preventDefault();
    var imgUrl = this.refs.imgUrl.getValue();
    if (imgUrl === this.cache.imgUrl) {
      return;
    }
    this.cache.imgUrl = imgUrl;
    this.predict(imgUrl);
    var href = this.makeHref('demo', {}, {img_url: imgUrl});
    window.location.href = href;
  },

  render: function() {
    var elems = this.state.bet_result.map(function(result) {
      var tag = result.tag;
      return (
        <tr>
          <td>{tag.tag_id}</td>
          <td><Link to="datasets" params={{dataType: 'all'}} query={{tag: tag.name}}>{tag.name}</Link></td>
          <td>{result.score}</td>
        </tr>
      );
    });

    var image;
    if (this.cache.imgUrl) {
      image = <img src={"/api/proxy?url=" + this.cache.imgUrl} />;
    }

    var time = 'loading...';
    if (this.state.time) {
      time = this.state.time + ' s';
    }

    return (
      <div className="dashboard demo">
        <Well>
          <h4>Enter an image url then predict tags.</h4>
          <form onSubmit={this.handleSubmit}>
            <Input type="url" name="img_url" ref="imgUrl" />
          </form>
        </Well>
        <Row className="result">
          <Col xs={6}>
            {image}
          </Col>
          <Col xs={6}>
            <Panel> Spend time: {time} </Panel>
            <Table striped bordered condensed hover>
              <thead>
                <tr>
                  <th>#</th>
                  <th>Tag</th>
                  <th>Score</th>
                </tr>
              </thead>
              <tbody>
                {elems}
              </tbody>
            </Table>
          </Col>
        </Row>
        <br />
        <br />
      </div>
    );
  }
});

var routes = (
  <Route handler={App} path="/">
    <DefaultRoute handler={Dashboard} />
    <Route name="dashboard" handler={Dashboard} />
    <Route name="demo" handler={DEMO} />
    <Route name="datasets" handler={Datasets} path="/ds/:dataType"/>
    <NotFoundRoute handler={Datasets} />
  </Route>
);

ReactRouter.run(routes, function (Handler, state) {
  React.render(<Handler />, document.body);
});
