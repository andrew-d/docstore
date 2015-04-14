jest.dontMock("../qpromise.js");

describe("queryablePromise", function() {
    var qpromise = require('../qpromise');

    // 'pit' == it with promise support
    pit("returns a promise when given one", function() {
        var p = new Promise(function(resolve, reject) {
            setTimeout(function() {
                resolve(1234);
            }, 1000);
        });

        var qp = qpromise(p);

        // Should have a .then() method
        expect(qp.then).toBeTruthy();

        jest.runAllTimers();
        return qp.then(function(result) {
            expect(result).toBe(1234);
        });
    });

    pit("reports the promise's status on success", function() {
        var p = qpromise(new Promise(function(resolve, reject) {
            setTimeout(function() {
                resolve(1234);
            }, 1000);
        }));

        expect(p.isFulfilled()).toBe(false);
        expect(p.isResolved()).toBe(false);
        expect(p.isRejected()).toBe(false);

        jest.runAllTimers();
        return p.then(function(val) {
            expect(p.isFulfilled()).toBe(true);
            expect(p.isResolved()).toBe(true);
            expect(p.isRejected()).toBe(false);
        });
    });

    pit("reports the promise's status on error", function() {
        var p = qpromise(new Promise(function(resolve, reject) {
            setTimeout(function() {
                reject("foo");
            }, 1000);
        }));

        expect(p.isFulfilled()).toBe(false);
        expect(p.isResolved()).toBe(false);
        expect(p.isRejected()).toBe(false);

        jest.runAllTimers();
        return p.then(null, function(val) {
            expect(p.isFulfilled()).toBe(true);
            expect(p.isResolved()).toBe(false);
            expect(p.isRejected()).toBe(true);
        });
    });

    pit("saves a promise's value after resolution", function() {
        var p = qpromise(new Promise(function(resolve, reject) {
            setTimeout(function() {
                resolve(1234);
            }, 1000);
        }));

        jest.runAllTimers();
        return p.then(function(val) {
            expect(p.value()).toBe(1234);
        });
    });

    pit("saves a promise's error", function() {
        var p = qpromise(new Promise(function(resolve, reject) {
            setTimeout(function() {
                reject("foo");
            }, 1000);
        }));

        jest.runAllTimers();
        return p.then(null, function(val) {
            expect(p.error()).toBe("foo");
        });
    });

    it("will throw when retrieving a non-resolved promise's value/error", function() {
        var p = qpromise(new Promise(function(resolve, reject) {
            setTimeout(function() {
                resolve(1234);
            }, 1000);
        }));

        expect(function() {
            p.value();
        }).toThrow("Promise is not resolved");

        expect(function() {
            p.error();
        }).toThrow("Promise is not rejected");
    });
});
