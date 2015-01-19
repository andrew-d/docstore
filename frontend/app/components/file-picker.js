import Ember from 'ember';
/* global mOxie */

export default Ember.Component.extend({
  tagName: 'button',
  classNameBindings: ['disabled'],

  // Configuration
  multiple: false,

  // Initializing m0xie can be a long operation - this ensures that we don't
  // let the button be clicked unless we're ready.
  disabled: true,

  didInsertElement: function() {
    var self = this;

    var fileInput = new mOxie.FileInput({
      browse_button: this.$().get(0),
      multiple:      this.get('multiple'),
    });

    fileInput.onchange = function() {
      self.sendAction('addFiles', fileInput.files);
    };

    fileInput.onready = function() {
      self.set('disabled', false);
    };

    fileInput.init();
  },
});
