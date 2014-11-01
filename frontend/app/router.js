import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.resource('tags', function() {
      this.route('show', { path: '/:tag_id' });
  });
});

export default Router;
