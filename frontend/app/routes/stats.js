import Ember from 'ember';

export default Ember.Route.extend({
  model: function() {
    return Ember.RSVP.hash({
      documents: this.store.find('document'),
      tags:      this.store.find('tag'),
      //files:     this.store.find('file'),
      files: [],
    });
  },

  setupController: function(controller, model) {
    // Set the various properties ('documents', 'tags', ...) on the controller
    // itself, as opposed to just under 'model'.
    controller.setProperties(model);
  }
});
