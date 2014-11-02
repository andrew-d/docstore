import DS from 'ember-data';

export default DS.Model.extend({
    ocr_data:   DS.attr('string'),
    created:    DS.attr('date'),
    image:      DS.belongsTo('image'),
    tags:       DS.hasMany('tag'),

    // Arbitrary other values
    metadata:   DS.attr('object'),
});
