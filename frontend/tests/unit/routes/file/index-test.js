import {
  moduleFor,
  test
} from 'ember-qunit';

moduleFor('route:file/index', 'FileIndexRoute', {
  // Specify the other units that are required for this test.
  // needs: ['controller:foo']
});

test('it exists', function() {
  var route = this.subject();
  ok(route);
});
