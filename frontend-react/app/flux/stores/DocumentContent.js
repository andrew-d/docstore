var ReactFlux = require('react-flux');

var DocumentConstants = require('../constants/Document'),
    TagConstants = require('../constants/Tag');

/**
 * This store is responsible for the content of all documents.  The state is an
 * Immutable map of Document ID ==> Document.
*/
var Store = ReactFlux.createStore({
  displayName: 'DocumentContentStore',

  getInitialState: function() {
    return {};
  },

  /**
   * This function loads any documents in the response into our store.  It
   * looks at response.body.{documents,document} and will load any found document.
   *
   * TODO: this should probably keep a LRU cache or something
   * TODO: currently, this just overwrites existing documents - we should
   *       probably do more intelligent merging (or at least warn for
   *       differences)
   */
  loadDocumentsFromResponse: function loadDocumentsFromResponse(resp) {
    if( !resp.body.documents && !resp.body.document ) {
      return;
    }

    var newState = this.state;

    if( resp.body.documents ) {
      resp.body.documents.forEach((t) => {
        newState = newState.set(t.id, t);
      });
    }

    if( resp.body.document ) {
      newState = newState.set(resp.body.document.id, resp.body.document);
    }

    if( newState !== this.state ) {
      this.setState(newState);
    }
  },

}, [
  // Listen to all server responses here.

  [TagConstants.LIST_SUCCESS, function onTagList(resp) {
    this.loadDocumentsFromResponse(resp);
  }],

  [TagConstants.GET_SUCCESS, function onTagGet(resp) {
    this.loadDocumentsFromResponse(resp);
  }],

  [DocumentConstants.LIST_SUCCESS, function onDocumentList(resp) {
    this.loadDocumentsFromResponse(resp);
  }],

  [DocumentConstants.GET_SUCCESS, function onDocumentGet(resp) {
    this.loadDocumentsFromResponse(resp);
  }],

  [DocumentConstants.CREATE_SUCCESS, function onDocumentCreate(resp) {
    this.loadDocumentsFromResponse(resp);
  }],
]);


module.exports = Store;
