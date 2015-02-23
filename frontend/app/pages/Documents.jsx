var React = require('react'),
    { Button, Input, Pager, PageItem } = require('react-bootstrap'),
    { Navigation, State } = require('react-router');

var DocumentActions = require('../flux').actions.Document;


const MIN_DOCUMENT_NAME = 5;


var Documents = React.createClass({
  mixins: [Navigation, State],

  getInitialState: function() {
    return {
      documentName: '',
    };
  },

  validationState: function() {
    var length = this.state.documentName.length;

    if( length === 0 ) {
      return undefined;
    } else if( length >= MIN_DOCUMENT_NAME ) {
      return 'success';
    } else {
      return 'error';
    }
  },

  handleDocNameChange: function() {
    this.setState({
      documentName: this.refs.docNameInput.getValue(),
    });
  },

  handleCreate: function() {
    var newName = this.state.documentName;

    if( newName.length < MIN_DOCUMENT_NAME ) {
      // TODO: show error
      return;
    }

    console.log("Creating document:", newName);
    DocumentActions.create(newName, (doc) => {
      console.log("Should transition to edit document:", doc);
    });
  },

  render: function() {
    var offset = this.getQuery().offset || 0,
        limit  = this.getQuery().limit || 20,
        nextUrl = this.makePath("documents", {}, {
          offset: offset + limit,       // TODO: overflow
          limit: limit,
        }),
        prevUrl = this.makePath("documents", {}, {
          offset: offset - limit,       // TODO: underflow
          limit: limit,
        });

    return (
      <div className="container-fluid">
        <div className="row">
          <div className="col-sm-9">
            <h3>Documents</h3>

            <Pager>
              <PageItem previous href={prevUrl}>&larr; Previous</PageItem>
              <PageItem next href={nextUrl}>Next &rarr;</PageItem>
            </Pager>
          </div>
          <div className="col-sm-3">
            <h3>Actions</h3>

            <h4>Create Document</h4>

            <Input
              type="text"
              value={this.state.documentName}
              placeholder="Enter document name"
              label="Document Name"
              ref="docNameInput"
              bsStyle={this.validationState()}
              onChange={this.handleDocNameChange}
            />

            <Button bsStyle="primary" onClick={this.handleCreate}>Create</Button>
          </div>
        </div>
      </div>
    );
  },
});


module.exports = Documents;
