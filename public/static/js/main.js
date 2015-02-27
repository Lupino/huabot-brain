var Navbar = ReactBootstrap.Navbar;
var Nav = ReactBootstrap.Nav;
var NavItem = ReactBootstrap.NavItem;
var MenuItem = ReactBootstrap.MenuItem;
var NavItemLink = ReactRouterBootstrap.NavItemLink;
var ButtonLink = ReactRouterBootstrap.ButtonLink;
var MenuItemLink = ReactRouterBootstrap.MenuItemLink;
var DropdownButton = ReactBootstrap.DropdownButton;
var ModalTrigger = ReactBootstrap.ModalTrigger;
var Modal = ReactBootstrap.Modal;
var Button = ReactBootstrap.Button;
var Input = ReactBootstrap.Input;
var Grid = ReactBootstrap.Grid;
var Row = ReactBootstrap.Row;
var Col = ReactBootstrap.Col;
var ListGroup = ReactBootstrap.ListGroup;
var ListGroupItem = ReactBootstrap.ListGroupItem;
var State = ReactRouter.State;
var Navigation = ReactRouter.Navigation;
var RouteHandler = ReactRouter.RouteHandler;
var Route = ReactRouter.Route;
var Router = ReactRouter.Router;
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
    jQuery("#waterfall").waterfall({
      selector: ".dataset",
    });
  },

  loadDatasets: function() {
    var self = this;
    var query = this.getQuery();
    var params = this.getParams();
    var max = query.max || '';
    var limit = query.limit || 50;
    var tag = query.tag || '';
    this.limit = limit;
    var dataType = params.dataType || 'all';
    jQuery.get('/api/datasets/?max=' + max + '&limit=' + limit + '&data_type=' + dataType + '&tag=' + tag, function(data) {
      self.setState(data);
    });
  },

  getInitialState: function() {
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
    var tag = this.getQuery().tag;
    if (this.cache.pathname !== pathname || this.cache.tag !== tag) {
      this.cache.pathname = pathname;
      this.cache.tag = tag;
      return true;
    }
    return false;
  },

  cleanDatasets: function() {
    this.datasets = [];
  },

  componentDidMount: function() {
    this.cache = this.cache || {};
    if (this.shouldLoadDatasets()) {
      this.loadDatasets();
    }
  },

  componentDidUpdate: function() {
    if (this.shouldLoadDatasets()) {
      if (this.shouldCleanDatasets()) {
        window.scroll(0, 0);
        this.cleanDatasets();
      }
      this.loadDatasets();
    } else {
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
              params={this.getParams()} to="datasets" query={query} block>加载更多...</ButtonLink>
        </div>
      );
    }
    if (this.datasets) {
      datasets = this.datasets.concat(datasets);
    }
    this.datasets = datasets;
    var elems = datasets.map(function(dataset) {
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
  getInitialState: function() {
    return {
      value: '',
      tags: []
    };
  },

  getHint: function(word) {
    var self = this;
    jQuery.get('/api/tags/hint?word=' + word, function(data) {
      self.setState(data);
    });
  },

  handleChange: function() {
    this.getHint(this.refs.tag.getValue());
    this.setState({
      value: this.refs.tag.getValue()
    });
  },

  handleSubmit: function(evt) {
    evt.preventDefault();
    this.props.onSubmit(this.state.value);
  },

  handleListClick: function(eventKey, href, target) {
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
    return {href: window.location.href};
  },

  handleSubmit: function(tag) {
    var query = this.getQuery();
    query.tag = tag;
    var href = this.makeHref('datasets', this.getParams(), query);
    window.location.href = href;
    this.setState({href: href});
  },
  render: function() {
    return (
      <div className="app-main">
        <Navbar fixedTop inverse fluid brand="Caffe Learn">
          <Nav right>
          </Nav>
          <SearchForm onSubmit={this.handleSubmit} />
        </Navbar>
        <Grid fluid>
          <Row>
            <Col sm={3} md={2} className="sidebar">
              <Nav>
                <NavItemLink to="datasets" params={{dataType: 'all'}}>所有数据</NavItemLink>
                <NavItemLink to="datasets" params={{dataType: 'train'}}>训练数据</NavItemLink>
                <NavItemLink to="datasets" params={{dataType: 'val'}}>验证数据</NavItemLink>
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

var routes = (
  <Route handler={App} path="/">
    <DefaultRoute handler={Datasets} />
    <Route name="datasets" handler={Datasets} path="/ds/:dataType"/>
    <NotFoundRoute handler={Datasets} />
  </Route>
);

ReactRouter.run(routes, function (Handler, state) {
  React.render(<Handler />, document.body);
});
