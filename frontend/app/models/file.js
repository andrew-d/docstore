import DS from 'ember-data';

export default DS.Model.extend({
  name:        DS.attr('string'),
  hash:        DS.attr('string'),
  size:        DS.attr('number'),
  created_at:  DS.attr('date'),

  // Type of file and associated properties.
  // NOTE: These properties are read-only
  type:        DS.attr('string'),
  properties:  DS.attr(),

  // TODO: can these be not async?
  tags:        DS.hasMany('tag', {async: true}),
  collections: DS.hasMany('collection', {async: true}),
});
