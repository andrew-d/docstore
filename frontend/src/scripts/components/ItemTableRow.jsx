var React = require('react'),
    Morearty = require('morearty'),
    moment = require('moment');


var ItemTableRow = React.createClass({
    mixins: [Morearty.Mixin],

    render: function() {
        var binding = this.getDefaultBinding(),
            niceTime = moment.utc(binding.get('created')).fromNow();

        return (
            <tr>
              <td>{binding.get('name')}</td>
              <td>{niceTime}</td>
              <td>{binding.get('files').count()}</td>
              <td>{binding.get('tags').toJS()}</td>
            </tr>
        );
    },
});

module.exports = ItemTableRow;
