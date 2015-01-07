var React = require('react'),
    component = require('omniscient');

var Home = component('Home', function() {
    return (
      <div className="jumbotron">
        <h1>docstore</h1>
        <p>A personal document store</p>
      </div>
    );
});

module.exports = Home;
