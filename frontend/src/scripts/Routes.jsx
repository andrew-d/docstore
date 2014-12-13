var React = require('react'),
    Router = require('react-router');

var Route = Router.Route,
    NotFoundRoute = Router.NotFoundRoute,
    DefaultRoute = Router.DefaultRoute,
    Link = Router.Link,
    RouteHandler = Router.RouteHandler;

var App = require('./App'),
    Items = require('./Items'),
    Home = require('./Home'),
    Stats = require('./Stats');


var routes = (
    <Route name="app" handler={App} path="/">
        <DefaultRoute name="index" handler={Home} />
        <Route name="items" handler={Items} path="/items" />
        <Route name="stats" handler={Stats} path="/stats" />
    </Route>
);

module.exports = routes;
