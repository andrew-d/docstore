/* global mOxie */

export function initialize() {
  // Set global configuration for mOxie
  mOxie.Env.swf_url = 'Moxie.min.swf';
  mOxie.Env.xap_url = 'Moxie.min.xap';
}

export default {
  name: 'moxie',
  initialize: initialize
};
