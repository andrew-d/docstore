import Ember from 'ember';
/* global moment */

export function momentCalendar(input) {
  return moment(input).calendar();
}

export default Ember.Handlebars.makeBoundHelper(momentCalendar);
