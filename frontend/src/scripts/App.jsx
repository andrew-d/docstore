var React = require('react'),
    component = require('omniscient'),
    RouteHandler = require('react-router').RouteHandler;

var Navbar = require('./Navbar').jsx;


var App = component('App', function() {
    return (
        <div>
          <Navbar />
          <div className="container">
            <RouteHandler {...this.props} />
          </div>
        </div>
    );
});

module.exports = App;
