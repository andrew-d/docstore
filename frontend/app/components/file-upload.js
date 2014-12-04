import Ember from 'ember';
/* global mOxie */

export default Ember.Component.extend({
  tagName: 'button',
  classNames: ['btn', 'btn-success'],

  // Configuration
  multiple: false,

  didInsertElement: function() {
    var self = this;

    var fileInput = new mOxie.FileInput({
      browse_button: this.$().get(0),
      multiple:      this.get('multiple'),
    });

    fileInput.onchange = function() {
      self.sendAction('addFiles', fileInput.files);
    };

    fileInput.init();
  },
});
