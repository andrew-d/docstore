import Ember from 'ember';
import RouteMixin from 'ember-cli-pagination/remote/route-mixin';

export default Ember.Route.extend(RouteMixin, {
  perPage: 25,

  model: function(params) {
    return this.findPaged('tag', params);
  }
});
