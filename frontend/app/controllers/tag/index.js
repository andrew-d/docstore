import Ember from 'ember';
import pagedArray from 'ember-cli-pagination/computed/paged-array';

export default Ember.Controller.extend({
  // Set up our query params
  queryParams: ["page", "perPage"],

  // Bind properties on the paged array to query params on this controller.
  pageBinding: "pagedDocuments.page",
  perPageBinding: "pagedDocuments.perPage",
  totalPagesBinding: "pagedDocuments.totalPages",

  // Default values for pagination.
  page: 1,
  perPage: 10,

  // TODO:
  haveMultiplePages: function() {
    return this.get('totalPages') > 1;
  }.property('totalPages'),

  // The paged content for this tag.
  pagedDocuments: function() {
    return pagedArray('model.documents', {perPage: this.get('perPage')});
  }.property('model', 'perPage'),
});
