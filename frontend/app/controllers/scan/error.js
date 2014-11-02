import Ember from 'ember';

export default Ember.Controller.extend({
    serverMessage: function() {
        return this.get('model.jqXHR.responseJSON').message;
    }.property('model.jqXHR.responseJSON'),
});
