import Ember from 'ember';
/* global moment */

export function momentFromNow(input) {
  return moment(input).fromNow();
}

export default Ember.Handlebars.makeBoundHelper(momentFromNow);
