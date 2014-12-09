var React = require('react'),
    Router = require('react-router'),
    routes = require('./scripts/Routes');

// We import this for the side-effect - it adds the .promise() method to the
// prototype of superagent.Request.
require('superagent-bluebird-promise');

Router.run(routes, function(Handler) {
    React.render(<Handler />, document.getElementById('application'));
});
