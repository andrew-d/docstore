import Ember from 'ember';

export default Ember.Route.extend({
  model: function(params) {
    var model = this.modelFor('file');    // Default behaviour

    // TODO: do we need this?  Seems like we always have a model
    if( !model ) {
      model = this.store.find('file', params.file_id);
    }

    return Ember.RSVP.hash({
      // The file for this route
      file: model,

      // All tags, for the tag select box
      tags: this.store.findAll('tag'),
    });
  },
});
