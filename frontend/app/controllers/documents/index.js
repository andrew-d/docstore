import Ember from 'ember';

export default Ember.Controller.extend({
    actions: {
        onScan: function() {
            // We kick off a scan by transitioning to the 'scan' route
            this.transitionToRoute("scan");
        },
    },
});
