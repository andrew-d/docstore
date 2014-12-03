import Ember from 'ember';

export default Ember.Controller.extend({
  fileSize: function() {
    return this.get('files').reduce(function(accum, item) {
      return accum + item.get('size');
    }, 0);
  }.property('files'),

  fileSizeStr: function() {
    return this.get('fileSize') + ' bytes';
  }.property('fileSize'),
});
