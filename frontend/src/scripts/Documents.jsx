var React = require('react'),
    Morearty = require('morearty'),
    State = require('react-router').State;

var Spinner = require('./components/Spinner');


var Documents = React.createClass({
    displayName: 'Documents',
    mixins: [Morearty.Mixin, State],

    componentDidMount: function() {
        // Reset the currently loaded documents when this component is mounted.
        this.getDefaultBinding().atomically()
            .set('documentsLoaded', false)
            .set('documents', [])
            .commit();

        // TODO: load new documents
        var query = this.getQuery(),
            currentPage = query.page || 1,
            perPage = query.per_page || 10;

        console.log("TODO: load documents");
        console.log("currentPage = " + currentPage);
        console.log("perPage = " + perPage);
    },

    render: function() {
        var b = this.getDefaultBinding();

        var documentsTable = b.get('documentsLoaded') ?
            <div>Documents Table</div> :
            <Spinner />;

        return (
            <div>
                <h1>Documents</h1>

                {documentsTable}
            </div>
        );
    },
});

module.exports = Documents;
