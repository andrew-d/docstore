var ReactFlux = require('react-flux'),
    request   = require('superagent'),
    promisify = require('superagent-promises');

var constants = require('../constants/Tag');

var Actions = ReactFlux.createActions({
  list: [constants.LIST, function() {
    return request.get('/api/documents')
                  .set('Accept', 'application/json')
                  .use(promisify)
                  .end();
  }],

  get: [constants.GET, function(id) {
    return request.get('/api/documents/' + id)
                  .set('Accept', 'application/json')
                  .end();
  }],

  create: [constants.CREATE, function(name, cb) {
    var p = request.post('/api/documents')
                   .send({name: name})
                   .set('Accept', 'application/json')
                   .use(promisify)
                   .end();

    // Callback when the document has been created.
    if( cb ) {
      p = p.then(function(resp) {
        cb(resp.body.document);
      });
    }

    return p;
  }],
});

module.exports = Actions;
