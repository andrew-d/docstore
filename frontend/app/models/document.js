import DS from 'ember-data';

export default DS.Model.extend({
  name: DS.attr('string'),
  created_at: DS.attr('date'),

  files: DS.hasMany('file'),
  tags: DS.hasMany('tag'),
  collection: DS.belongsTo('collection'),
});
