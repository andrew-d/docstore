function when(promise, spec) {
    // If the promise does not have our query function (see: qpromise.js),
    // then we can't use this helper.
    if( !promise.isResolved ) {
        throw new Error("Promise is not queryable");
    }

    // Depending on the state...
    if( promise.isResolved() ) {
        // Done
        if( spec.done ) {
            return spec.done(promise.value());
        }
    } else if( promise.isRejected() ) {
        // Failed
        if( spec.failed ) {
            return spec.failed(promise.error());
        }
    } else {
        // Pending
        if( spec.pending ) {
            return spec.pending();
        }
    }

    // Return null by default - this indicates to React.js that we don't want
    // to render anything.
    return null;
}

module.exports = when;
