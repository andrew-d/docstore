import Ember from 'ember';
/* global Humanize */

export function humanFileSize(input) {
  return Humanize.fileSize(input);
}

export default Ember.Handlebars.makeBoundHelper(humanFileSize);
