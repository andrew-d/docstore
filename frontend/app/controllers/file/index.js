import Ember from 'ember';
import { getOrCreateTag } from '../../util/tag-util';

export default Ember.Controller.extend({
  // New tag to add, bound to select box
  tagToAdd: null,

  // Property that filters out tags already in the file.
  // TODO: we don't currently use this
  nonexistantTags: function() {
    return this.model.tags.reject((item) => {
      return this.model.file.get('tags').contains(item);
    });
  }.property('model.tags', 'model.file.tags'),

  actions: {
    removeTag: function(tag) {
      this.model.file.get('tags').removeObject(tag);
      this.model.file.save();
    },

    // Triggered whenever ember-selectize wants to create a nonexistant tag.
    newTag: function(tag) {
      getOrCreateTag(this.store, tag);
    },

    // Triggered when the add tag form is submitted.
    addTag: function() {
      this.store.find('tag', {name: this.get('tagToAdd')})
          .then((tagObj) => this.model.file.get('tags').pushObject(tagObj));
    }
  },
});
