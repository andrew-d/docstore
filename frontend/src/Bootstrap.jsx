var React = require('react'),
    Handler = require('react-router').Handler,
    Ctx = require('./scripts/Ctx'),
    Routes = require('./scripts/Routes');

// Problem: how to integrate Morearty and react-router?
//
// 1. We need to render Morearty like so:
//      var Bootstrap = Ctx.bootstrap(Handler);
//      React.render(<Bootstrap />, document.getElementById('application'));
//
// 2. However, to get the router working, we need to do this:
//      Router.run(Routes, function(Handler) {
//          React.render(<Handler />, document.getElementById('application'));
//      });
//
// What to do?

// TODO: fixme
React.render(
    <div>Temporary</div>,
    document.getElementById('application')
);
