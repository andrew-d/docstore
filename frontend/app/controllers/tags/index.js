import Ember from 'ember';
import Paginated from '../../mixins/paginated';

export default Ember.Controller.extend(Paginated, {
  total: function(){
    return this.store.metadataFor('tag').total;
  }.property('model')
});
