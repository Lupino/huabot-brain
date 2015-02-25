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
var Grid = ReactBootstrap.Grid;
var Row = ReactBootstrap.Row;
var Col = ReactBootstrap.Col;
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
  mixins: [ReactRouter.State],

  waterfall: function() {
    jQuery("#waterfall").waterfall({
      selector: ".img",
    });
  },

  loadDatasets: function() {
    var self = this;
    var query = this.getQuery();
    var params = this.getParams();
    var max = query.max || '';
    var limit = query.limit || 50;
    this.limit = limit;
    var dataType = params.dataType || 'all';
    jQuery.get('/api/datasets/?max=' + max + '&limit=' + limit + '&data_type=' + dataType, function(data) {
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
    if (this.cache.pathname !== pathname) {
      this.cache.pathname = pathname;
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
      loadMore = (
        <div className="load-more">
          <ButtonLink bsStyle="info" bsSize="large" params={this.getParams()} to="datasets" query={{max: datasets[datasets.length - 1].dataset_id, limit: this.limit}}>加载更多...</ButtonLink>
        </div>
      );
    }
    if (this.datasets) {
      datasets = this.datasets.concat(datasets);
    }
    this.datasets = datasets;
    var elems = datasets.map(function(dataset) {
      return (
        <ModalTrigger modal={<Dataset data={dataset} title={dataset.tag.name} />}>
          <div className="img" data-id={dataset.dataset_id}>
            <div className="dataset">
              <img src={"/upload/" + dataset.file.key} width={192} height={192/dataset.file.width * dataset.file.height} />
              <div className="tag">{dataset.tag.name}</div>
            </div>
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

var App = React.createClass({
  render: function() {
    return (
      <div className="app-main">
        <Navbar fixedTop inverse fluid brand="Caffe Learn">
          <Nav>
          </Nav>
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
