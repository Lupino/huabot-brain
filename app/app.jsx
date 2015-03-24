var App = React.createClass({
  contextTypes: {
    router: React.PropTypes.func.isRequired
  },

  mixins: [OverlayMixin],

  getInitialState () {
    var {router} = this.context;
    var tag = router.getCurrentQuery().tag || '';
    return {tag: tag, isModalOpen: false};
  },

  handleSubmit(tag) {
    var {router} = this.context;
    var query = router.getCurrentQuery();
    var params = router.getCurrentParams();
    params.dataType = params.dataType || 'all';
    query.tag = tag;
    delete query.max;
    var href = router.makeHref('datasets', params, query);
    window.location.href = href;
    this.setState({tag: tag || ''});
  },

  handleToggle (evt) {
    this.setState({
      isModalOpen: !this.state.isModalOpen
    });
  },

 // This is called by the `OverlayMixin` when this component
  // is mounted or updated and the return value is appended to the body.
  renderOverlay () {
    if (!this.state.isModalOpen) {
      return <span />;
    }

    return <NewDataset onRequestHide={this.handleToggle} />;
  },
  render () {
    return (
      <div className="app-main">
        <Navbar fixedTop inverse fluid brand="Huabot Brain">
          <Nav right>
            <NavItemLink to="dashboard">Dashboard</NavItemLink>
            <NavItem onClick={this.handleToggle} href={null}><Glyphicon glyph="plus" /></NavItem>
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
