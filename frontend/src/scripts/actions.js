var Reflux = require('reflux'),
    data = require('./data');


var rootUrl = 'http://localhost:8080';


var actions = Reflux.createActions({
    'loadItems': ['completed', 'failed'],
});


actions.loadItems.listen(function(options) {
    actions.loadItems.promise(
        request
            .get(rootUrl + '/api/items')
            .query({
                page:     options.currentPage || 1,
                per_page: options.perPage || 20,
            })
            .promise()
    );
});

actions.loadItems.completed.listen(function(res) {
    // TODO: update data
});

// TODO: handle failure


module.exports = actions;
