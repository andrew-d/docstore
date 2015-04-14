var React = require('react'),
    Router = require('react-router');

var { Route, DefaultRoute } = Router;

// Require routes
var About = require('./pages/About'),
    App = require('./pages/App'),
    Documents = require('./pages/Documents'),
    Home = require('./pages/Home');


var Routes = (
  <Route handler={App} path="/">
    {/* Introduction page */}
    <DefaultRoute name="index" handler={Home} />

    {/* Display documents (paginated) */}
    <Route name="documents" handler={Documents} />

    {/* About page */}
    <Route name="about" handler={About} />
  </Route>
);


module.exports = Routes;
