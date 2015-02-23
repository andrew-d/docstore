var ReactFlux = require('react-flux');

var Constants = ReactFlux.createConstants([
  'LIST',
  'GET',
  'CREATE',
], 'DOCUMENT');

module.exports = Constants;
