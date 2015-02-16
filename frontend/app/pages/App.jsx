var React = require('react'),
    RouteHandler = require('react-router').RouteHandler;

var ReactBootstrap = require('react-bootstrap'),
    Nav = ReactBootstrap.Nav,
    Navbar = ReactBootstrap.Navbar;

var ReactRouterBootstrap = require('react-router-bootstrap'),
    NavItemLink = ReactRouterBootstrap.NavItemLink;


var App = React.createClass({
  render: function() {
    return (
      <div className='page-wrapper'>
        <Navbar fluid={true} staticTop={true} brand='Docstore'>
          <Nav>
            <NavItemLink to='index'>
              Home
            </NavItemLink>
            <NavItemLink to='about'>
              About
            </NavItemLink>
          </Nav>
        </Navbar>

        <div className='content-container'>
          <RouteHandler />
        </div>
      </div>
    );
  },
});


module.exports = App;
