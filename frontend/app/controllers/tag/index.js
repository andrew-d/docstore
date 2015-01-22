import Ember from 'ember';

// TODO: this should be paginated - we currently show all files on one page.
export default Ember.Controller.extend({
  queryParams: ['display'],
  display: 'grid',

  isGrid: function() {
    return this.get('display') === 'grid';
  }.property('display'),

  isList: function() {
    return this.get('display') === 'list';
  }.property('display'),
});
