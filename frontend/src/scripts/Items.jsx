var React = require('react'),
    component = require('omniscient'),
    State = require('react-router').State,
    request = require('superagent');

var Spinner = require('./components/Spinner').jsx,
    ItemTable = require('./components/ItemTable').jsx,
    Pagination = require('./components/Pagination');


var Items = component('Items', State, function(props) {
    var itemsTable = props.cursor.get('itemsLoaded') ?
        <ItemTable binding={props.cursor.sub('items')} /> :
        <Spinner />;

    return (
        <div>
            <h1>Items</h1>

            {itemsTable}

            {/* TODO: make this really paginate */}
            <Pagination currPage={1} totalPages={30} />
        </div>
    );
});

module.exports = Items;
