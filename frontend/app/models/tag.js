import DS from 'ember-data';

export default DS.Model.extend({
  name:  DS.attr('string'),

  // TODO: can this be not async?
  files: DS.hasMany('file', {async: true}),
});
