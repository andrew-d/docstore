import Ember from 'ember';

export default Ember.Controller.extend({
  queryParams: ['file_id'],

  // Query param values
  file_id: 0,

  // Current file is a propery of the ID
  file: function() {
    return null;
  }.property('file_id'),
});
