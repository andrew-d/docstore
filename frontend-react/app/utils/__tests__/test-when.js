jest.dontMock("../when.js");
jest.dontMock("../qpromise.js");

describe("when", function() {
    var when = require("../when"),
        qpromise = require('../qpromise'),
        successPromise,
        failPromise;

    beforeEach(function() {
        successPromise = qpromise(new Promise(function(resolve, reject) {
            setTimeout(function() {
                resolve("success");
            }, 1000);
        }));

        failPromise = qpromise(new Promise(function(resolve, reject) {
            setTimeout(function() {
                reject("fail");
            }, 1000);
        }));
    });

    it("will throw if given a non-queryable promise", function() {
        expect(function() {
            when(new Promise(function() {}), {});
        }).toThrow("Promise is not queryable");
    });

    it("will handle non-resolved promises", function() {
        var pendingMock = jest.genMockFunction();

        when(successPromise, {
            pending: pendingMock,
        });

        expect(pendingMock.mock.calls.length).toBe(1);
        expect(pendingMock.mock.calls[0]).toEqual([]);
    });

    pit("will handle resolved promises", function() {
        jest.runAllTimers();
        return successPromise.then(function() {
            var doneMock = jest.genMockFunction();

            var ret = when(successPromise, {
                done: doneMock.mockReturnValueOnce(1234),
            });

            expect(doneMock.mock.calls.length).toBe(1);
            expect(doneMock.mock.calls[0][0]).toBe("success");
            expect(ret).toBe(1234);
        });
    });

    pit("will handle rejected promises", function() {
        jest.runAllTimers();
        return failPromise.then(null, function() {
            var failedMock = jest.genMockFunction();

            var ret = when(failPromise, {
                failed: failedMock.mockReturnValueOnce(5678),
            });

            expect(failedMock.mock.calls.length).toBe(1);
            expect(failedMock.mock.calls[0][0]).toBe("fail");
            expect(ret).toBe(5678);
        });
    });

    pit("will return null for methods that are not defined", function() {
        // Pending
        var ret = when(successPromise, {});
        expect(ret).toBe(null);

        jest.runAllTimers();
        return successPromise.then(function() {
            // Success
            var ret = when(successPromise, {});
            expect(ret).toBe(null);
        });
    });

    pit("will return null for methods that are not defined (2)", function() {
        // Pending
        var ret = when(failPromise, {});
        expect(ret).toBe(null);

        jest.runAllTimers();
        return failPromise.then(null, function() {
            // Fail
            var ret = when(failPromise, {});
            expect(ret).toBe(null);
        });
    });
});
