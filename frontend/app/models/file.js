import DS from 'ember-data';

export default DS.Model.extend({
  name:        DS.attr('string'),
  size:        DS.attr('number'),
  created_at:  DS.attr('date'),

  // TODO: can these be not async?
  tags:        DS.hasMany('tag', {async: true}),
  collections: DS.hasMany('collection', {async: true}),
});
