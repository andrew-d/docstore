import Ember from 'ember';
import { request as rawRequest, defineFixture as rawDefineFixture } from 'ic-ajax';
import config from '../config/environment';

// This util prefixes requests in development with a different hostname.  This
// lets us run our backend on a different port, and thus still make use of the
// LiveReload features that 'ember serve' provides by default.

var prefix;
if( config.environment === 'development' ) {
    prefix = 'http://localhost:8888';
} else {
    prefix = '';
}

export function request(/*arguments*/) {
    var args = Array.prototype.slice.call(arguments);

    if( typeof args[0] === 'string' ) {
        args[0] = prefix + args[0];
    } else if( args[0].hasOwnProperty('url') && typeof args[0].url === 'string' ) {
        args[0].url = prefix + args[0].url;
    } else {
        throw new Ember.Error("unknown invocation of ic-ajax");
    }

    return rawRequest.apply(null, args);
}

export default request;

// defineFixture() should also prefix the URL
var defineFixture;
if( config.environment === 'development' ) {
    defineFixture = function(url, fixture) {
        return rawDefineFixture(prefix + url, fixture);
    };
} else {
    defineFixture = function() {};
}

export var defineFixture = defineFixture;
