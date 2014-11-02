import DS from 'ember-data';

export default DS.Model.extend({
    hash:       DS.attr('string'),
    path:       DS.attr('string'),
    'document': DS.belongsTo('document'),
    created:    DS.attr('date'),
});
