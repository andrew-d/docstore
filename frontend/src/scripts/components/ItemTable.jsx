var React = require('react'),
    Morearty = require('morearty'),
    map = require('lodash-node/modern/collections/map');

var ItemTableRow = require('./ItemTableRow');


var ItemTable = React.createClass({
    mixins: [Morearty.Mixin],

    render: function() {
        var binding = this.getDefaultBinding();

        var tableRows = map(binding.get().toJS(), function(item, index) {
            return <ItemTableRow key={item.id}
                                 binding={binding.sub(index)} />
        });

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
    },
});

module.exports = ItemTable;
