import Ember from 'ember';

export function initialize(container, application) {
  application.deferReadiness();

  Ember.$.getJSON("/api/capabilities", function(json) {
    var Capabilities = Ember.Object.extend(json);

    // Inject capabilities into all controllers.
    application.register('capabilities:main', Capabilities);
    application.inject('controller', 'capabilities', 'capabilities:main');

    application.advanceReadiness();
  });
}

export default {
  name: 'capabilities',
  initialize: initialize
};
