import Ember from 'ember';

export default Ember.ArrayController.extend({
  // Set up our query params
  queryParams: ["page", "perPage"],

  // Bind properties on the paged array to query params on this controller.
  pageBinding: "content.page",
  perPageBinding: "content.perPage",
  totalPagesBinding: "content.totalPages",

  // Whether or not we have more than one page of items
  haveMultiplePages: function() {
    return this.get('totalPages') > 1;
  }.property('totalPages'),
});
