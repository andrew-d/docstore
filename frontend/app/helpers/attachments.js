import Ember from 'ember';
/* global mOxie */

// Note: helper functions from:
//    https://github.com/scribu/ember-moxie-demo/blob/gh-pages/app.js

/**
 * Load an attachment into memory.
 *
 * @param Blob file - The file to load
 * @return RSVP.Promise
 */
export function loadAttachment(file) {
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
export function uploadAttachment(file, data, url) {
  return new Ember.RSVP.Promise(function(resolve, reject) {
    var req = Ember.$.post(url, {
      filename: file.name,
      data: data
    });

    function successHandler(response) {
      resolve(response);
    }

    function errorHandler(xhr) {
      reject(xhr.responseText);
    }

    req.then(successHandler, errorHandler);
  });
}
