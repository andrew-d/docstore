import Ember from 'ember';
/* global mOxie, _ */

/**
 * Load a file into memory.
 *
 * @param Blob file - The file to load
 * @return RSVP.Promise
 */
function loadFile(file) {
  return new Ember.RSVP.Promise(function(resolve, reject) {
    var reader = new mOxie.FileReader();

    reader.onloadend = function() {
      resolve(reader.result);
    };

    reader.onerror = function() {
      reject(reader.error);
    };

    reader.readAsBinaryString(file);
  });
}

/**
 * Send a file to the server.
 *
 * @param Blob file - The file info
 * @param String data - The binary file content
 * @return RSVP.Promise
 */
function uploadFile(url, file, data, additional) {
  return new Ember.RSVP.Promise(function(resolve, reject) {
    var opts = additional || {},
        params = _.extend({}, opts, {
          filename: file.name,
          data:     data,
        });

    var req = Ember.$.post(url, params);

    var successHandler = function successHandler(response) {
      resolve(response);
    };

    var errorHandler = function errorHandler(xhr) {
      reject(xhr.responseText);
    };

    req.then(successHandler, errorHandler);
  });
}

export default Ember.Controller.extend({
  // Files from the file picker.
  files: [],

  // Bindings for other things
  tags: [],
  collection: null,

  // Values set when we upload
  uploading: false,
  completed: 0,
  errors: [],

  actions: {
    // Add some new files
    addFiles: function(files) {
      var existing = this.get('files');

      files.forEach(function(file) {
        if( !existing.findBy('name', file.name) ) {
          existing.pushObject(file);
        }
      });
    },

    // Remove a single file
    removeFile: function(file) {
      var existing = this.get('files');
      this.set('files', existing.rejectBy('name', file.name));
    },

    // Submit all files to the server, and then clear the list of
    // files and show any errors.
    submit: function() {
      var errors = [],
          files = this.get('files');

      // TODO: put this somewhere else
      var UPLOAD_URL = '/api/files/upload';

      if( Ember.isEmpty(files) ) {
        this.notify.alert("No files selected to upload!");
        return;
      }

      this.set('uploading', true);

      // Map all files to a promise that is resolved or rejected when
      // that file is uploaded.
      var promises = files.map((file) => {
        return loadFile(file)
          .then((data) => {
            // TODO: tags
            // TODO: collection
            return uploadFile(UPLOAD_URL, file, data);
          })
          .then(() => {
            this.incrementProperty('completed');
          }, (error) => {
            errors.push({ file: file.name, reason: error });
          });
      });

      // When we've uploaded (or failed) all files, then we set the appropriate
      // variables on our controller.
      Ember.RSVP.all(promises).then(() => {
        if( errors.length ) {
          this.set('errors', errors);
        }

        this.set('uploading', false);
        this.set('files', []);
      });
    },

    // Clear the list of files.
    reset: function() {
      this.set('files', []);
      this.set('completed', 0);
    },

    // Action that is sent when an new tag is created in the select box.
    newTag: function(tag) {
      // TODO: dedupe code from here and controller file/index.js
      this.store
        .find('tag', {name: tag})
        .then((tags) => {
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
          var record = this.store.createRecord('tag', {
            name: tag,
          });
          return record.save();
        })
      .then((tagObj) => {
        this.set('model.tags', this.store.findAll('tag'));
        this.get('tags').pushObject(tagObj);
        this.notify.info("Created new tag: " + tagObj.get('name'));
      });
    }
  },
});
