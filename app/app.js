var App = React.createClass({
  mixins: [State, Navigation],

  getInitialState: function() {
    var tag = this.getQuery().tag || '';
    return {tag: tag};
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

