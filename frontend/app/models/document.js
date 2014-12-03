import DS from 'ember-data';

export default DS.Model.extend({
    name: DS.attr('string'),
    created: DS.attr('date'),
    files: DS.hasMany('file'),
    meta: DS.attr(),
    tags: DS.hasMany('tag'),
});
