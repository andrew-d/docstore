var React = require('react'),
    Immutable = require('immutable'),
    Morearty = require('morearty'),
    State = require('react-router').State,
    request = require('superagent');

var Spinner = require('./components/Spinner');


var Items = React.createClass({
    displayName: 'Items',
    mixins: [Morearty.Mixin, State],

    componentDidMount: function() {
        var self = this;

        // Reset the currently loaded items when this component is mounted.
        this.getDefaultBinding().atomically()
            .set('itemsLoaded', false)
            .set('items', [])
            .commit();

        var query = this.getQuery(),
            currentPage = query.page || 1,
            perPage = query.per_page || 10;

        // TODO: move this to a flux action - e.g. Reflux
        request
            .get('http://localhost:8080' + '/api/items')
            .query({page: currentPage, per_page: perPage})
            .promise()
            .then(function(res) {
                self.getDefaultBinding().atomically()
                    .set('itemsLoaded', true)
                    .set('items', Immutable.fromJS(res.body))
                    .commit();
            })
            .catch(function(e) {
                // TODO: handle errors
                console.log(e);
            });
    },

    render: function() {
        var b = this.getDefaultBinding();

        // TODO: render real table of items
        var itemsTable = b.get('itemsLoaded') ?
            <div>Items Table</div> :
            <Spinner />;

        return (
            <div>
                <h1>Items</h1>

                {itemsTable}
            </div>
        );
    },
});

module.exports = Items;
