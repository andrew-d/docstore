var Morearty = require('morearty');

var Ctx = Morearty.createContext({
    // Global store for objects from the server.
    documents: [],
    tags: [],
    files: [],
});

module.exports = Ctx;
