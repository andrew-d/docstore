import DS from 'ember-data';

export default DS.Model.extend({
    size: DS.attr('number'),
    name: DS.attr('string'),
    created: DS.attr('date'),
    'document': DS.belongsTo('document'),
});
