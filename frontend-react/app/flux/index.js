module.exports = {
  constants: {
    Document: require('./constants/Document'),
    Tag:      require('./constants/Tag'),
  },

  actions: {
    Document: require('./actions/Document'),
    Tag:      require('./actions/Tag'),
  },

  stores: {
    TagContent: require('./stores/TagContent'),
    TagDocuments: require('./stores/TagDocuments'),
    DocumentContent: require('./stores/DocumentContent'),
  },
};
