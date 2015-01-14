import DS from 'ember-data';

export default DS.Model.extend({
  name:  DS.attr('string'),
  files: DS.hasMany('file'),
});
