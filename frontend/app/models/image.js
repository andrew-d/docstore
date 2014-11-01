import DS from 'ember-data';

export default DS.Model.extend({
    hash:       DS.attr('string'),
    path:       DS.attr('string'),
    receipt:    DS.belongsTo('receipt'),
    created:    DS.attr('date'),
});
