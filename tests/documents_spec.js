var expect = require('chai').expect,
    request = require('superagent');

var config = require('./config');
require('./spec_helper.js');

describe('documents', function() {
  it('has no documents by default', function(done) {
    request
      .get(config.URL+'/documents')
      .end(function(res) {
        expect(res.ok).to.be.ok;
        expect(res.body).to.deep.equal({
          documents: [],
          tags: [],
        });

        done();
      });
  });

  describe('basic creation', function() {
    var createdId;

    beforeEach(function(done) {
      request
        .post(config.URL + '/documents')
        .send({name: 'document name'})
        .end(function(res) {
          expect(res.ok).to.be.ok;
          expect(res.body).to.have.property('document');
          expect(res.body.document.name).to.equal('document name');

          createdId = res.body.document.id;
          done();
        });
    });

    afterEach(function(done) {
      request
        .del(config.URL + '/documents/' + createdId)
        .end(function(res) {
          done();
        });
    });

    it('should allow fetching the document directly', function(done) {
      request
        .get(config.URL + '/documents/' + createdId)
        .end(function(res) {
          expect(res.ok).to.be.ok;
          expect(res.body.document.id).to.equal(createdId);
          expect(res.body.document.name).to.equal('document name');
          expect(res.body.document.created_at).to.be.a('number');
          expect(res.body.document.tags).to.deep.equal([]);
          expect(res.body.tags).to.deep.equal([]);

          done();
        });
    });

    it('should display the document when fetching all documents', function(done) {
      request
        .get(config.URL + '/documents')
        .end(function(res) {
          expect(res.ok).to.be.ok;
          expect(res.body.documents).to.have.length(1);

          var doc = res.body.documents[0]
          expect(doc.id).to.equal(createdId);
          expect(doc.name).to.equal('document name');
          expect(doc.created_at).to.be.a('number');
          expect(doc.tags).to.deep.equal([]);

          done();
        });
    });

    it('should disallow creating a document with a too-short name', function(done) {
      request
        .post(config.URL + '/documents')
        .send({name: '11'})
        .end(function(res) {
          expect(res.status).to.equal(400);
          done();
        });
    });
  });
});
