import Ember from 'ember';
import request from '../../utils/ajax';

export default Ember.Route.extend({
    model: function() {
        return request({
            method: 'POST',
            url: '/api/images/scan'
        });
    },
});
