export function getOrCreateTag(store, tagName) {
  return store
    .find('tag', {name: tagName})
    .then((tags) => {
      // The returned value should be an array with exactly 1 element.
      if( tags.get('length') !== 1 ) {
        throw new Error("successful tag lookup should return 1 object");
      }

      return tags.objectAt(0);
    }, (reason) => {
      if( reason.status !== 404 ) {
        throw reason;
      }

      // No tag by this name - create it.
      var record = store.createRecord('tag', {
        name: tagName,
      });
      return record.save();
    });
}
