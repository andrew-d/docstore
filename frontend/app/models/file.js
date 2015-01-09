import DS from 'ember-data';

export default DS.Model.extend({
  name:       DS.attr('string'),
  size:       DS.attr('number'),
  created_at: DS.attr('date'),
  item:       DS.belongsTo('item'),
});
