import DS from 'ember-data';

export default DS.Model.extend({
  name: DS.attr('string'),

  // TODO: can this be not async?
  files: DS.hasMany('file', {async: true}),

  // Reflexive relations!  These need to be async - otherwise, we'd need
  // to serialize the entire collection graph for each request.
  children: DS.hasMany('collection', {inverse: 'parent', async: true}),
  parent:   DS.belongsTo('collection', {inverse: 'children', async: true}),
});
