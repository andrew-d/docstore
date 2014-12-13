var Morearty = require('morearty');

var Ctx = Morearty.createContext({
    // Items list page
    items: [],
    itemsLoaded: false,

    // Tags list page
    tags: [],

    // Files list page
    files: [],
});

module.exports = Ctx;
