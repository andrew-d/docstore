import Ember from 'ember';

export function joinTags(tags, connector) {
  return tags.reduce(function(accum, tag, idx) {
    // On the first iteration, don't add a connector.
    if( idx === 0 ) {
      return tag.get('name');
    }

    return accum + connector + tag.get('name');
  }, '');
}

export default Ember.Handlebars.makeBoundHelper(joinTags);
