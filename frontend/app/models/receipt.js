import DS from 'ember-data';

export default DS.Model.extend({
    amount:     DS.attr('number'),
    ocr_data:   DS.attr('string'),
    created:    DS.attr('date'),
    image:      DS.belongsTo('image'),
    tags:       DS.hasMany('tag'),
});
