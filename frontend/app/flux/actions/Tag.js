var ReactFlux = require('react-flux'),
    request = require('superagent-bluebird-promise');

var constants = require('../constants/Tag');

var Actions = ReactFlux.createActions({
  fetch: [constants.FETCH, function() {
    return request.get('/api/tags')
                  .set('Accept', 'application/json')
                  .promise();
  }],

  create: [constants.CREATE, function(name) {
    return request.post('/api/tags')
                  .send({name: name})
                  .set('Accept', 'application/json')
                  .promise();
  }],
});

module.exports = Actions;
