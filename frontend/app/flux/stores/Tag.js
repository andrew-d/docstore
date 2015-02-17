var ReactFlux = require('react-flux');
var constants = require('../constants/Tag');


var Store = ReactFlux.createStore({
  displayName: 'TagStore',

  getInitialState: function(){
    return {};
  },
}, [

  [constants.FETCH_SUCCESS, function onFetch(tags) {
    console.log("Tags:", tags);
  }],

  [constants.CREATE_SUCCESS, function onCreate(tag) {
    console.log("Created tag:", tag);
  }],

]);


module.exports = Store;
