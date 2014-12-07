var React = require('react'),
    Morearty = require('morearty');

var Home = React.createClass({
    mixins: [Morearty.Mixin],

    render: function() {
        return (
          <div className="jumbotron">
            <h1>docstore</h1>
            <p>A personal document store</p>
          </div>
        );
    },
});

module.exports = Home;
