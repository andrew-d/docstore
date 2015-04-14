import Ember from 'ember';
/* global moment */

export function momentFromNow(params) {
  return moment(params[0]).fromNow();
}

export default Ember.HTMLBars.makeBoundHelper(momentFromNow);
