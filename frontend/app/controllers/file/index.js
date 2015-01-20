import Ember from 'ember';

export default Ember.Controller.extend({
  newTag: null,

  actions: {
    removeTag: function(tag) {
      this.model.get('tags').removeObject(tag);
      this.model.save();
    },

    addTag: function() {
      var self = this,
          newTag = this.get('newTag');

      self.store
        .find('tag', {name: newTag})
        .then(function(tags) {
          // The returned value should be an array with exactly 1 element.
          if( tags.get('length') !== 1 ) {
            throw new Error("successful tag lookup should return 1 object");
          }

          return tags.objectAt(0);
        }, function(reason) {
          if( reason.status !== 404 ) {
            throw reason;
          }

          // No tag by this name - create it.
          var record = self.store.createRecord('tag', {
            name: newTag,
          });
          return record.save();
        })
      .then(function(tag) {
        self.model.get('tags').addObject(tag);
        self.model.save();
        self.set('newTag', null);
      });
    },
  },
});
