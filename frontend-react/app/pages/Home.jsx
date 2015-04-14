var React = require('react');

var tagActions = require('../flux/actions/Tag');


var Home = React.createClass({
  render: function() {
    return (
      <div>
        <h1>React Webpack Starter</h1>

        <p>This is the home page.</p>

        <button onClick={this.handleList}>
          List Tags
        </button>
        <button onClick={this.handleCreate}>
          Create Tag
        </button>
      </div>
    );
  },

  handleList: function() {
    tagActions.list();
  },

  handleCreate: function() {
    tagActions.create("Test Tag");
  },
});


module.exports = Home;
