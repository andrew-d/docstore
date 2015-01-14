import DS from 'ember-data';

export default DS.Model.extend({
  name: DS.attr('string'),
  files: DS.hasMany('file'),

  // Reflexive relations!
  children: DS.hasMany('collection', {inverse: 'parent'}),
  parent:   DS.belongsTo('collection', {inverse: 'children'}),
});
