import Ember from 'ember';

export default Ember.Component.extend({
  // Which attribute we're sorting on.
  sortAttribute: 'created_at',
  sortAscending: true,

  sortedDocuments: function() {
    var sorted = this.get('documents').sortBy(this.get('sortAttribute')).toArray();

    if( !this.get('sortAscending') ) {
      sorted = sorted.reverse();
    }

    return sorted;
  }.property('documents', 'sortAttribute', 'sortAscending'),


  // Properties for the sort classes
  nameSortClass: function() {
    return this._sortPropertyFor('name');
  }.property('sortAttribute', 'sortAscending'),

  createdAtSortClass: function() {
    return this._sortPropertyFor('created_at');
  }.property('sortAttribute', 'sortAscending'),

  // Helper function for sort properties
  _sortPropertyFor: function(attr) {
    if( this.get('sortAttribute') !== attr ) {
      return 'fa-sort';
    }

    if( this.get('sortAscending') ) {
      return 'fa-sort-desc';
    } else {
      return 'fa-sort-asc';
    }
  },

  actions: {
    setSort: function(attr) {
      // If we're already sorting by this, toggle the ascending.
      if( this.get('sortAttribute') === attr ) {
        this.toggleProperty('sortAscending');
        return false;
      }

      // Set this as the sort property, reset 'ascending' flag
      this.set('sortAttribute', attr);
      this.set('sortAscending', true);
      return false;
    },
  },
});
