import Ember from 'ember';

export default Ember.ArrayController.extend({
  // Set up our query params
  queryParams: ["page", "perPage"],

  // Bind properties on the paged array to query params on this controller.
  pageBinding: "content.page",
  perPageBinding: "content.perPage",
  totalPagesBinding: "content.totalPages",

  // Default values for pagination.
  page: 1,
  perPage: 10,

  // Whether or not we have more than one page of items
  haveMultiplePages: function() {
    return this.get('totalPages') > 1;
  }.property('totalPages'),

  // Active tab helpers
  tabName: 'upload',

  isUploadActive: function() {
      return this.get('tabName') === 'upload';
  }.property('tabName'),

  isScanActive: function() {
      return this.get('tabName') === 'scan';
  }.property('tabName'),

  // Values for creating a new document
  newName: null,
  newTags: null,
  newFile: null,
  scannerName: null,

  // Files we're uploading
  files: [],

  actions: {
    // Switch to the given tab
    switchTab: function(name) {
      if( ['upload', 'scan'].indexOf(name) === -1 ) {
        throw new Error("Invalid tab: " + name);
      }

      this.set('tabName', name);
    },

    // Add new file to our list of files.  Sent from 'file-upload' component.
    addFiles: function(files) {
      var existing = this.get('files');

      // Add new files to our array.
      files.forEach(function(file) {
        if( !existing.findBy('name', file.name) ) {
          existing.pushObject(file);
        }
      });
    },

    // Remove a file from our list of files
    removeFile: function(file) {
      var existing = this.get('files');
      this.set('files', existing.rejectBy('name', file.name));
    },

    // TODO: submit new files
    submit: function() {
      console.log("Would submit new");
    },

    // TODO: scan new files
    scan: function() {
      console.log("Would scan new files");
    },

    reset: function() {
      this.set('newName', null);
      this.set('newTags', null);
      this.set('newFile', null);
      this.set('scannerName', null);
      this.set('files', []);
    },
  },
});
