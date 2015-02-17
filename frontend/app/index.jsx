var React = require('react'),
    Router = require('react-router'),
    routes = require('./Routes');

// Require the global stylesheet.
require('./styles/index.scss');

// Require flux components.
require('./flux');

if( process.env.NODE_ENV !== "production" ) {
    // Dev tool support
    window.React = React;
}

Router.run(routes, Router.HistoryLocation, function(Handler) {
  React.render(<Handler />, document.body);
});
