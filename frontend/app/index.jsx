var Marty = require('marty'),
    React = require('react'),
    Router = require('react-router'),
    routes = require('./Routes');

// Require the global stylesheet.
require('./styles/index.scss');

if( process.env.NODE_ENV !== "production" ) {
    // Dev tool support
    window.React = React;
    window.Marty = Marty;
}

Router.run(routes, Router.HistoryLocation, function(Handler) {
  React.render(<Handler />, document.body);
});
