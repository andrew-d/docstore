import Ember from 'ember';
import Paginated from '../../mixins/paginated';

export default Ember.Controller.extend(Paginated, {
  queryParams: ['display'],
  display: 'grid',

  isGrid: function() {
    return this.get('display') === 'grid';
  }.property('display'),

  isList: function() {
    return this.get('display') === 'list';
  }.property('display'),

  total: function(){
    return this.store.metadataFor('item').total;
  }.property('model')
});
