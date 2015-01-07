var immstruct = require('immstruct');

var data = immstruct({
    // Items list page
    items: [],
    itemsLoaded: false,

    // Tags list page
    tags: [],

    // Files list page
    files: [],
});

module.exports = data;
