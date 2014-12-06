var React = require('react'),
    Morearty = require('morearty'),
    RouteHandler = require('react-router').RouteHandler;

var App = React.createClass({
    displayName: 'App',
    mixins: [Morearty.Mixin],

    render: function() {
        // Always render the child for now.
        return <RouteHandler />;
    },
});

module.exports = App;
