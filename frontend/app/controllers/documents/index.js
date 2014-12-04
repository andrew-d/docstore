import Ember from 'ember';

export default Ember.ArrayController.extend({
  // Set up our query params
  queryParams: ["page", "perPage"],

  // Bind properties on the paged array to query params on this controller.
  pageBinding: "content.page",
  perPageBinding: "content.perPage",
  totalPagesBinding: "content.totalPages",

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

  // Available scanners
  // TODO: where should we set this?
  availableScanners: [
    {id: 1, label: "foo"},
    {id: 2, label: "bar"},
  ],

  actions: {
    switchTab: function(name) {
      if( ['upload', 'scan'].indexOf(name) === -1 ) {
        throw new Error("Invalid tab: " + name);
      }

      this.set('tabName', name);
    },

    // Sent from 'file-upload' component
    addFiles: function(files) {
      var existing = this.get('files');

      // Add new files to our array.
      files.forEach(function(file) {
        if (!existing.findBy('name', file.name)) {
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
  },
});
