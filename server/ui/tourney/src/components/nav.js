import React, { Component } from 'react';
import { Link } from 'react-router-dom'

export default class NavBar extends Component {
  render() {
    return (
      <div>
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
        <NavBreadcrumbs navPath={this.props.location.pathname}/>
      </div>
    ); 
  }
}

class NavBreadcrumbs extends Component {
  render() {
    let path = this.props.navPath;
    if (path.length > 0 && path[0] === '/') {
      path = path.substring(1, path.length);
    }
    if (path.length > 0 && path[path.length-1] === '/') {
      path = path.substring(0, path.length-1);
    }

    var crumbItems = [];
    let pathTerms = path.split('/');
    for (let i=0; i < pathTerms.length; i++) {
      if (pathTerms[i] !== '') {
        let link = '/' + pathTerms.slice(0, i+1).join('/');
        let term = pathTerms[i].charAt(0).toUpperCase() + pathTerms[i].slice(1);
        if (i < pathTerms.length - 1) {
          crumbItems.push(<li key={term}><Link to={link}>{term}</Link></li>);
        } else {
        crumbItems.push(<li className="active" key={term}>{term}</li>);
        }
      }
    }

    return (
      <ul className="breadcrumb">
        { crumbItems }
      </ul>
    );
  }
}
