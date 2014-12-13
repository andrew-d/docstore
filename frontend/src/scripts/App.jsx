var React = require('react'),
    Morearty = require('morearty'),
    RouteHandler = require('react-router').RouteHandler;

var Ctx = require('./Ctx'),
    Navbar = require('./Navbar');

var App = React.createClass({
    displayName: 'App',

    componentDidMount: function() {
        Ctx.init(this);

        // Debugging
        if( process.env.NODE_ENV !== 'production' ) window.Ctx = Ctx;
    },

    render: function() {
        return React.withContext({ morearty: Ctx }, function() {
            return (
                <div>
                  <Navbar />
                  <div className="container">
                    <RouteHandler binding={Ctx.getBinding()} />
                  </div>
                </div>
            );
        });
    },
});

module.exports = App;
