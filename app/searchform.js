var SearchForm = React.createClass({
  contextTypes: {
    router: React.PropTypes.func.isRequired
  },

  getInitialState () {
    var {router} = this.context;
    var query = router.getCurrentQuery();
    this.cache = this.cache || {};
    this.cache.changed = true;
    this.cache.path = router.getCurrentPath();

    return {
      value: query.tag || '',
      tags: []
    };
  },

  getHint (word) {
    var self = this;
    jQuery.get('/api/tags/hint?word=' + word, data => self.setState(data));
  },

  shouldCleanHint () {
    var path = this.context.router.getCurrentPath();
    if (this.cache.path === path) {
      return false;
    }
    this.cache.path = path;
    return true;

  },

  cleanHint () {
    this.setState({tags: []});
  },

  componentDidUpdate () {
    if (this.cache.changed) {
      if (this.shouldCleanHint()) {
        this.cache.changed = false;
        this.cleanHint();
      }
      return;
    }
    var query = this.context.router.getCurrentQuery();
    if (query.tag && query.tag !== this.state.value) {
      this.setState({
        value: query.tag
      });
    } else if (!query.tag && this.state.value) {
      this.setState({value: ''});
    }
    this.cache.path = this.context.router.getCurrentPath();
  },

  handleChange () {
    this.cache.changed = true;
    this.getHint(this.refs.tag.getValue());
    this.setState({
      value: this.refs.tag.getValue()
    });
  },

  handleSubmit (evt) {
    evt.preventDefault();
    this.cache.changed = false;
    this.props.onSubmit(this.state.value);
  },

  handleListClick (evt) {
    var target = evt.target.innerText.trim();
    this.cache.changed = false;
    this.props.onSubmit(target);
    this.setState({tags: [], value: target});
  },

  render () {
    var list = this.state.tags.map(tag => <ListGroupItem> {tag.name} </ListGroupItem>);
    return (
      <form className="navbar-form navbar-right" onSubmit={this.handleSubmit}>
        <Input type="text"
          name="tag"
          value={this.state.value}
          ref="tag"
          placeholder="Search..."
          onChange={this.handleChange} />

        <div className="hint" onClick={this.handleListClick}>
          <ListGroup>
            {list}
          </ListGroup>
        </div>
      </form>
    );
  }
});
