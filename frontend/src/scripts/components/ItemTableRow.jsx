var React = require('react'),
    component = require('omniscient'),
    moment = require('moment');


var ItemTableRow = component('ItemTableRow', function(props) {
    var c = props.cursor,
        niceTime = moment.utc(c.get('created')).fromNow();

    return (
        <tr>
          <td>{c.get('name')}</td>
          <td>{niceTime}</td>
          <td>{c.get('files').count()}</td>
          <td>{c.get('tags').toJS()}</td>
        </tr>
    );
});

module.exports = ItemTableRow;
