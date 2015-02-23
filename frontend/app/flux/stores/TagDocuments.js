var Immutable = require('immutable'),
    ReactFlux = require('react-flux');


/**
 * This store is responsible for keeping track of which documents beling to
 * which tags.  The state is a mapping of Tag ID ==> Immutable List of
 * Document IDs.
 */
var Store = ReactFlux.createStore({
  displayName: 'TagDocumentsStore',

  getInitialState: function() {
    return {};
  },

}, [
  [TagConstants.CREATE_SUCCESS, function onTagCreate(resp) {
    // When we've created a tag, there are no documents attached to it.
    var newState = this.state.set(resp.tag.id, Immutable.List());
    this.setState(newState);
  }],

  [TagConstants.GET_SUCCESS, function onTagGet(resp) {
    // TODO: Extract documents from this tag
  }],

  [TagConstants.LIST_SUCCESS, function onTagList(resp) {
    // TODO: Extract documents from all tags
  }],
]);


module.exports = Store;
