import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('stats');
  this.resource('documents', function() {
    this.resource('document', { path: '/:document_id' }, function() { });
  });
  this.resource('tags', function() {
    this.resource('tag', { path: '/:tag_id' }, function() { });
  });
});

export default Router;
