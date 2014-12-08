var Morearty = require('morearty');

var Ctx = Morearty.createContext({
    // Documents list page
    documents: [],
    documentsPage: 1,
    documentsPerPage: 10,

    // Tags list page
    tags: [],

    // Files list page
    files: [],
});

module.exports = Ctx;
