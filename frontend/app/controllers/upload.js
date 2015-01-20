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
      var self = this,
          errors = [],
          files = self.get('files');

      // TODO: put this somewhere else
      var UPLOAD_URL = '/api/files/upload';

      if( Ember.isEmpty(files) ) {
        self.notify.alert("No files selected to upload!");
        return;
      }

      self.set('uploading', true);

      // Map all files to a promise that is resolved or rejected when
      // that file is uploaded.
      var promises = files.map(function(file) {
        return loadFile(file)
          .then(function(data) {
            // TODO: tags
            // TODO: collection
            return uploadFile(UPLOAD_URL, file, data);
          })
          .then(function success() {
            self.incrementProperty('completed');
          }, function error(error) {
            errors.push({ file: file.name, reason: error });
          });
      });

      // When we've uploaded (or failed) all files, then we set the appropriate
      // variables on our controller.
      Ember.RSVP.all(promises).then(function() {
        if( errors.length ) {
          self.set('errors', errors);
        }

        self.set('uploading', false);
        self.set('files', []);
      });
    },

    // Clear the list of files.
    reset: function() {
      this.set('files', []);
      this.set('completed', 0);
    },
  },
});
