var React = require('react');

var tagActions = require('../flux/actions/Tag');


var Home = React.createClass({
  render: function() {
    return (
      <div>
        <h1>React Webpack Starter</h1>

        <p>This is the home page.</p>

        <button onClick={this.handleFetch}>
          Fetch Tags
        </button>
        <button onClick={this.handleCreate}>
          Create Tag
        </button>
      </div>
    );
  },

  handleFetch: function() {
    tagActions.fetch();
  },

  handleCreate: function() {
    tagActions.create("Test Tag");
  },
});


module.exports = Home;
