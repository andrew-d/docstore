var React = require('react'),
    Router = require('react-router');

var App = require('./scripts/App'),
    data = require('./scripts/data'),
    { requestAnimationFrame, cancelAnimationFrame } = require('./scripts/requestAnimationFrame');


// We import this for the side-effect - it adds the .promise() method to the
// prototype of superagent.Request.
require('superagent-bluebird-promise');

var rerender = function rerender(structure, el) {
    var Handler, state;

    var render = function render(h, s) {
        // TODO: this is b0rked
        if (h) Handler = h;
        if (s) state = s;

        React.render(<Handler cursor={structure.cursor()}
                              routePath={state.path}
                              statics={state} />, el);
    };

    // Only rerender on requestAnimationFrame.
    var queuedChange = false;
    structure.on('swap', function structureSwapped() {
        if( queuedChange ) return;
        queuedChange = true;

        requestAnimationFrame(function forceRender() {
            queuedChange = false;
            // TODO: arguments?
            render();
        });
    });

    return render;
};


var routes = require('./scripts/Routes');
Router.run(routes, rerender(data, document.getElementById('application')));
