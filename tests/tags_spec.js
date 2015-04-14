var expect = require('chai').expect,
    request = require('superagent');

var config = require('./config');
require('./spec_helper.js');


describe('tags', function() {
  it('has no tags by default', function(done) {
    request
      .get(config.URL+'/tags')
      .end(function(res) {
        expect(res.ok).to.be.ok;
        expect(res.body).to.deep.equal({
          tags: [],
        });

        done();
      });
  });

  describe('basic creation', function() {
    var createdId;

    beforeEach(function(done) {
      request
        .post(config.URL + '/tags')
        .send({name: 'foo'})
        .end(function(res) {
          expect(res.ok).to.be.ok;
          expect(res.body).to.have.property('tag');
          expect(res.body.tag.name).to.equal('foo');

          createdId = res.body.tag.id;
          done();
        });
    });

    afterEach(function(done) {
      request
        .del(config.URL + '/tags/' + createdId)
        .end(function(res) {
          done();
        });
    });

    it('should allow fetching the tag directly', function(done) {
      request
        .get(config.URL + '/tags/' + createdId)
        .end(function(res) {
          expect(res.ok).to.be.ok;
          expect(res.body).to.deep.equal({
            tag: {
              id: createdId,
              name: "foo",
              documents: [],
            },
          });

          done();
        });
    });

    it('should display the tag when fetching all tags', function(done) {
      request
        .get(config.URL + '/tags')
        .end(function(res) {
          expect(res.ok).to.be.ok;
          expect(res.body).to.deep.equal({
            tags: [
              {
                id: createdId,
                name: "foo",
                documents: [],
              },
            ],
          });

          done();
        });
    });

    it('should disallow creating a tag with the same name', function(done) {
      request
        .post(config.URL + '/tags')
        .send({name: 'foo'})
        .end(function(res) {
          expect(res.ok).to.not.be.ok;
          done();
        });
    });

    it('should allow fetching only specific tags', function(done) {
      // Create a second tag
      request
        .post(config.URL + '/tags')
        .send({name: 'bar'})
        .end(function(res) {
          expect(res.ok).to.be.ok;

          var secondId = res.body.tag.id;

          // Get only this tag.
          request
            .get(config.URL + '/tags')
            .query({ 'ids[]': secondId })
            .end(function(res) {
              expect(res.ok).to.be.ok;
              expect(res.body).to.deep.equal({
                tags: [
                  {
                    id: secondId,
                    name: 'bar',
                    documents: [],
                  },
                ],
              });

              // Remove this tag.
              request
                .del(config.URL + '/tags/' + secondId)
                .end(function(res) {
                  done();
                });
            });
        });
    });
  });
});
