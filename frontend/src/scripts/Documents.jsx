var React = require('react'),
    Morearty = require('morearty');

var Spinner = require('./components/Spinner');


var Documents = React.createClass({
    displayName: 'Documents',
    mixins: [Morearty.Mixin],

    componentDidMount: function() {
        console.log("Did mount");
    },

    render: function() {
        return (
            <div>
                <h1>Documents</h1>

                <Spinner />
            </div>
        );
    },
});

module.exports = Documents;
