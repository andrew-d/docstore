import DS from 'ember-data';

export default DS.Model.extend({
  name: DS.attr('string'),
  hash: DS.attr('string'),
  'document': DS.belongsTo('document'),
});
