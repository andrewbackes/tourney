import React, { Component } from 'react';
import { Link } from 'react-router-dom'

export default class NavBar extends Component {
  render() {
    return (
      <nav className="navbar navbar-inverse">
        <div className="container-fluid">
          <div className="navbar-header">
          <button type="button" className="navbar-toggle" data-toggle="collapse" data-target="#myNavbar">
            <span className="icon-bar"></span>
            <span className="icon-bar"></span>
            <span className="icon-bar"></span>                        
          </button>
          <a className="navbar-brand" href="">Tourney</a>
          </div>
          <div className="collapse navbar-collapse" id="myNavbar">
          <ul className="nav navbar-nav">
            <li className={this.props.location.pathname.startsWith("/tournaments") ? 'active' : ''}><Link to='/tournaments'>Tournaments</Link></li>
            <li className={this.props.location.pathname.startsWith("/engines") ? 'active' : ''}><Link to='/engines'>Engines</Link></li>
            <li className={this.props.location.pathname.startsWith("/books") ? 'active' : ''}><Link to='/books'>Books</Link></li>
            <li className={this.props.location.pathname.startsWith("/workers") ? 'active' : ''}><Link to='/workers'>Workers</Link></li>
          </ul>
          <ul className="nav navbar-nav navbar-right">
            <li><Link to='/signup'><span className="glyphicon glyphicon-user"></span> Sign Up</Link></li>
            <li><Link to='/login'><span className="glyphicon glyphicon-log-in"></span> Login</Link></li>
          </ul>
          </div>
        </div>
      </nav>
    ); 
  }
}