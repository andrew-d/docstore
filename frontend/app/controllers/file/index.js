import Ember from 'ember';

export default Ember.Controller.extend({
  newTag: null,

  actions: {
    removeTag: function(tag) {
      console.log("Would remove tag: " + tag.get('name'));
    },

    addTag: function() {
      var self = this,
          newTag = this.get('newTag');

      var addTag = function(tag) {
        self.model.get('tags').pushObject(tag);
        self.model.save();
        self.set('newTag', null);
      };

      self.store.find('tag', {name: newTag})
                .then(addTag, function(reason) {
                  if( reason.status !== 404 ) {
                    throw reason;
                  }

                  var record = self.store.createRecord('tag', {
                    name: newTag,
                  });
                  record.save().then(addTag);
                });
    },
  },
});
