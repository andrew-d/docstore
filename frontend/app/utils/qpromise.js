// Create a promise that allows querying its state and final values.
var queryablePromise = function queryablePromise(promise) {
    // Don't create a wrapper for promises that can already be queried.
    if (promise.isResolved) return promise;

    var isResolved = false,
        isRejected = false,
        promiseResult,
        promiseError;

    // Observe the promise, saving the fulfillment in a closure scope.
    var result = promise.then(
        function(v) {
            isResolved = true;
            promiseResult = v;
            return v;
        },
        function(e) {
            isRejected = true;
            promiseError = e;
            throw e;
        }
    );

    result.isFulfilled = function() { return isResolved || isRejected; };
    result.isResolved = function() { return isResolved; }
    result.isRejected = function() { return isRejected; }

    result.value = function() {
        if( !isResolved ) {
            throw new Error("Promise is not resolved");
        }

        return promiseResult;
    };
    result.error = function() {
        if( !isRejected ) {
            throw new Error("Promise is not rejected");
        }

        return promiseError;
    };

    return result;
};


module.exports = queryablePromise;
