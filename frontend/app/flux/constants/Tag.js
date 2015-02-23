var ReactFlux = require('react-flux');

var Constants = ReactFlux.createConstants([
  'LIST',
  'GET',
  'CREATE',
], 'TAG');

module.exports = Constants;
