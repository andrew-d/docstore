var extend = require('lodash-node/modern/object/assign');


/* Application-wide configuration */
var config = {
    // Configuration goes here
    foo: 'bar',
};


if( 'production' !== process.env.NODE_ENV ) {
  extend(config, {
    baseUrl: 'http://localhost:4201',
  });
} else {
  extend(config, {
    baseUrl: '',
  });
}


module.exports = config;
