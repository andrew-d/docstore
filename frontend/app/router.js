import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.resource('tags', function() {
      this.route('show', { path: '/:tag_id' });
  });
  this.resource('documents',  function() {
      this.route('show', { path: '/:document_id' });
      this.resource('scan', { path: '/scan' }, function() { });
  });
});

export default Router;
