import Ember from 'ember';
/* global Humanize */

export function humanIntComma(input) {
  return Humanize.intComma(input);
}

export default Ember.Handlebars.makeBoundHelper(humanIntComma);
