var Morearty = require('morearty');

var Ctx = Morearty.createContext({
    // Documents list page
    documents: [],
    documentsLoaded: false,

    // Tags list page
    tags: [],

    // Files list page
    files: [],
});

module.exports = Ctx;
