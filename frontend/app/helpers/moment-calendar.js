import Ember from 'ember';
/* global moment */

export function momentCalendar(params) {
  return moment(params[0]).calendar();
}

export default Ember.HTMLBars.makeBoundHelper(momentCalendar);
