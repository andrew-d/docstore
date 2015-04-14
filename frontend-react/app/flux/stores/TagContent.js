var ReactFlux = require('react-flux');

var DocumentConstants = require('../constants/Document'),
    TagConstants = require('../constants/Tag');

/**
 * This store is responsible for the content of all tags.  The state is an
 * Immutable map of Tag ID ==> Tag.
*/
var Store = ReactFlux.createStore({
  displayName: 'TagContentStore',

  getInitialState: function() {
    return {};
  },

  /**
   * This function loads any tags in the response into our store.  It looks
   * at response.body.{tags,tag} and will load any found tag.
   *
   * TODO: this should probably keep a LRU cache or something
   * TODO: currently, this just overwrites existing tags - we should probably
   *       do more intelligent merging (or at least warn for differences)
   */
  loadTagsFromResponse: function loadTagsFromResponse(resp) {
    if( !resp.body.tags && !resp.body.tag ) {
      return;
    }

    var newState = this.state;

    if( resp.body.tags ) {
      resp.body.tags.forEach((t) => {
        newState = newState.set(t.id, t);
      });
    }

    if( resp.body.tag ) {
      newState = newState.set(resp.body.tag.id, resp.body.tag);
    }

    if( newState !== this.state ) {
      this.setState(newState);
    }
  },

}, [
  // Listen to all server responses here.

  [TagConstants.LIST_SUCCESS, function onTagList(resp) {
    this.loadTagsFromResponse(resp);
  }],

  [TagConstants.GET_SUCCESS, function onTagGet(resp) {
    this.loadTagsFromResponse(resp);
  }],

  [TagConstants.CREATE_SUCCESS, function onTagCreate(resp) {
    this.loadTagsFromResponse(resp);
  }],

  [DocumentConstants.LIST_SUCCESS, function onDocumentList(resp) {
    this.loadTagsFromResponse(resp);
  }],

  [DocumentConstants.GET_SUCCESS, function onDocumentGet(resp) {
    this.loadTagsFromResponse(resp);
  }],
]);


module.exports = Store;
