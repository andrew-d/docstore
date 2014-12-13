var React = require('react'),
    Morearty = require('morearty');


var Stats = React.createClass({
    mixins: [Morearty.Mixin],

    renderScanners: function() {
        var b = this.getDefaultBinding();

        // TODO: fix this
        if( true ) {
            return <p><strong>No scanners found.</strong></p>;
        }

        var scanners = b.get('scanners').map(function(s) {
            return (
              <tr>
                <td>TKTK Name</td>
                <td>TKTK Vendor</td>
                <td>TKTK Model</td>
              </tr>
            );
        });

        return (
          <table className="table table-condensed">
            <thead>
              <th>Name</th>
              <th>Vendor</th>
              <th>Model</th>
            </thead>
            <tbody>
                {scanners}
            </tbody>
          </table>
        );
    },

    render: function() {
        var b = this.getDefaultBinding();
        var maybeScanners = this.renderScanners();

        return (
            <div className="row">
              <div className="col-xs-12">
                <h2>Statistics</h2>

                There are {b.get('items').count()} items(s) in the store.<br/>
                There are {b.get('files').count()} file(s) saved, totalling TKTK.<br/>
                There are {b.get('tags').count()} tag(s) in the store.

                <h3>Scanners</h3>
                {maybeScanners}
              </div>
            </div>
        );
    },
});

module.exports = Stats;
