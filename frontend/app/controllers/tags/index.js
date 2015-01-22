import Ember from 'ember';
import Paginated from '../../mixins/paginated';

export default Ember.ArrayController.extend(Paginated, {
  sortProperties: ['name'],
  sortAscending: true,

  total: function(){
    return this.store.metadataFor('tag').total;
  }.property('model')
});
