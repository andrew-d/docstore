var React = require('react'),
    Router = require('react-router');

var Route = Router.Route,
    NotFoundRoute = Router.NotFoundRoute,
    DefaultRoute = Router.DefaultRoute,
    Link = Router.Link,
    RouteHandler = Router.RouteHandler;

var App = require('./App'),
    Home = require('./Home');


var routes = (
    <Route handler={App} path="/">
        <DefaultRoute handler={Home} />
    </Route>
);

module.exports = routes;
