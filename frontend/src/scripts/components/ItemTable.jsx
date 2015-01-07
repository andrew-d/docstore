var React = require('react'),
    component = require('omniscient'),
    map = require('lodash-node/modern/collections/map');

var ItemTableRow = require('./ItemTableRow').jsx;


var ItemTable = component('ItemTable', function(props) {
    var tableRows = props.cursor.map(function(item, index) {
        return <ItemTableRow key={item.get('id')}
                             cursor={props.cursor.cursor(index)} />
    }).toArray();

    return (
      <table className="table table-striped table-bordered table-condensed">
        <thead>
          <th>Name</th>
          <th>Created</th>
          <th># of Files</th>
          <th>Tags</th>
        </thead>
        <tbody>
          {tableRows}
        </tbody>
      </table>
    );
});

module.exports = ItemTable;
