var React = require('react'),
    Router = require('react-router');

var { Route, DefaultRoute } = Router;

// Require routes
var About = require('./pages/About'),
    App = require('./pages/App'),
    Home = require('./pages/Home');


var Routes = (
  <Route handler={App} path="/">
    {/* Introduction page */}
    <DefaultRoute name="index" handler={Home} />

    {/* About page */}
    <Route name="about" handler={About} />
  </Route>
);


module.exports = Routes;
