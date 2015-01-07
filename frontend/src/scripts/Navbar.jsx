var React = require('react'),
    component = require('omniscient'),
    Router = require('react-router'),
    Link = Router.Link,
    Navigation = Router.Navigation,
    State = Router.State;


// Returns true whenever the route changes.
var ShouldComponentUpdateMixin = {
    shouldComponentUpdate: function(newProps, newState) {
        return (this.props.to       !== newProps.to       ||
                this.props.children !== newProps.children ||
                component.shouldComponentUpdate(newProps, newState)
               );
    },
};


var NavbarLink = component('NavbarLink', [Navigation, State, ShouldComponentUpdateMixin], function(props) {
    var isActive = this.isActive(props.to, props.params, props.query || null);
    var className = isActive ? "active" : "";

    return (
        <li className={className}>
          <Link {...props} children={undefined}>
            {props.children}
          </Link>
        </li>
    );
}).jsx;


var Navbar = component('Navbar', [State, ShouldComponentUpdateMixin], function() {
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
              <NavbarLink to="items">Items</NavbarLink>
              <NavbarLink to="stats">Stats</NavbarLink>
            </ul>
          </div>
        </div>
      </nav>
    );
});

module.exports = Navbar;
