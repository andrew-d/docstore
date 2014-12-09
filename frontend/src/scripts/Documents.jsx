var React = require('react'),
    Immutable = require('immutable'),
    Morearty = require('morearty'),
    State = require('react-router').State,
    request = require('superagent');

var Spinner = require('./components/Spinner');


var Documents = React.createClass({
    displayName: 'Documents',
    mixins: [Morearty.Mixin, State],

    componentDidMount: function() {
        var self = this;

        // Reset the currently loaded documents when this component is mounted.
        this.getDefaultBinding().atomically()
            .set('documentsLoaded', false)
            .set('documents', [])
            .commit();

        var query = this.getQuery(),
            currentPage = query.page || 1,
            perPage = query.per_page || 10;

        // TODO: move this to a flux action - e.g. Reflux
        request
            .get('http://localhost:8080' + '/api/documents')
            .query({page: currentPage, per_page: perPage})
            .promise()
            .then(function(res) {
                self.getDefaultBinding().atomically()
                    .set('documentsLoaded', true)
                    .set('documents', Immutable.fromJS(res.body))
                    .commit();
            })
            .catch(function(e) {
                // TODO: handle errors
                console.log(e);
            });
    },

    render: function() {
        var b = this.getDefaultBinding();

        // TODO: render real table of documents
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
