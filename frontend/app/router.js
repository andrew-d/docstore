import Ember from "ember";
import config from "./config/environment";

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.resource("files", function() {
    this.resource("file", {
      path: ":file_id"
    }, function() {});
  });

  this.route("upload");

  this.resource("tags", function() {
    this.resource("tag", {
      path: ":tag_id"
    }, function() {});
  });
});

export default Router;