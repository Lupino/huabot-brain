var SearchForm = React.createClass({
  mixins: [State],

  getInitialState: function() {
    var query = this.getQuery();
    this.cache = this.cache || {};
    this.cache.changed = true;
    this.cache.path = this.getPath();

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
