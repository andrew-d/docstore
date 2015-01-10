import Ember from 'ember';
import PaginationBase from '../pagination-base';
/* global Holder */

export default PaginationBase.extend({
  init: function(){
    this._super('item');
  },

  // This allows us to use Holder.js placeholders in the images.
  renderTemplate: function() {
    Ember.run.schedule('afterRender', null, function () { Holder.run(); });
    this._super.apply(this, arguments);
  },
});
