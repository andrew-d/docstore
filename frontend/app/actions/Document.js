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

  create: [constants.CREATE, function(name) {
    return request.post('/api/tags')
                  .send({name: name})
                  .set('Accept', 'application/json')
                  .use(promisify)
                  .end();
  }],
});

module.exports = Actions;
