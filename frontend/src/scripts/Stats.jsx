var React = require('react'),
    component = require('omniscient');


var Stats = component('Stats', function(props) {
    var c = props.cursor;

    var renderScanners = function() {
        // TODO: fix this
        if( true ) {
            return <p><strong>No scanners found.</strong></p>;
        }

        var scanners = c.get('scanners').map(function(s) {
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
    };

    var maybeScanners = renderScanners();

    return (
        <div className="row">
          <div className="col-xs-12">
            <h2>Statistics</h2>

            There are {c.get('items').count()} items(s) in the store.<br/>
            There are {c.get('files').count()} file(s) saved, totalling TKTK.<br/>
            There are {c.get('tags').count()} tag(s) in the store.

            <h3>Scanners</h3>
            {maybeScanners}
          </div>
        </div>
    );
});

module.exports = Stats;
