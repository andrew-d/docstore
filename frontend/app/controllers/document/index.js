import Ember from 'ember';

export default Ember.Controller.extend({
  // Active tab panel
  activeTab: "tag",


  fileSize: function() {
    return this.get('model.files').reduce(function(accum, item) {
      return accum + item.get('size');
    }, 0);
  }.property('model.files'),

  fileSizeStr: function() {
    return this.get('fileSize') + ' bytes';
  }.property('fileSize'),

  actions: {
    deleteTag: function(tag) {
      this.model.get('tags').removeObject(tag);
      this.model.save();
    },
  },
});
