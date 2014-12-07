var React = require('react'),
    Router = require('react-router'),
    Link = Router.Link,
    Navigation = Router.Navigation,
    State = Router.State;


var NavbarLink = React.createClass({
    mixins: [Navigation, State],

    _isActive: function() {
        return this.isActive(this.props.to, this.props.params, this.props.query || null);
    },

    render: function() {
        var className = this._isActive() ? "active" : "";

        return (
            <li className={className}>
              <Link {...this.props} children={undefined}>
                {this.props.children}
              </Link>
            </li>
        );
    },
});


var Navbar = React.createClass({
    render: function() {
        return (
          <nav className="navbar navbar-default navbar-static-top" role="navigation">
            <div className="container">
              <div className="navbar-header">
                <button type="button" className="navbar-toggle collapsed"
                        data-toggle="collapse" data-target="#navbar"
                        aria-expanded="false" aria-controls="navbar">
                  <span className="sr-only">Toggle navigation</span>
                  <span className="icon-bar"></span>
                  <span className="icon-bar"></span>
                  <span className="icon-bar"></span>
                </button>
                <Link to="app" className="navbar-brand">
                  Docstore
                </Link>
              </div>
              <div id="navbar" className="navbar-collapse collapse">
                <ul className="nav navbar-nav">
                  <NavbarLink to="index">Home</NavbarLink>
                </ul>
                <ul className="nav navbar-nav">
                  <NavbarLink to="stats">Stats</NavbarLink>
                </ul>
              </div>
            </div>
          </nav>
        );
    },
});

module.exports = Navbar;
