import Ember from 'ember';

export default Ember.ArrayController.extend({
    actions: {
        createTag: function() {
            var tagName = this.get('newTag');
            if( !tagName ) { return false; }
            if( !tagName.trim() ) { return; }

            // Create the new Tag model.
            var tag = this.store.createRecord('tag', {
                name: tagName,
            });

            // Clear text field
            this.set('newTag', '');

            // Save model
            tag.save();
        },
    },
});
