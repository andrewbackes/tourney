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
              <li className='active'><Link to='/tournaments'>Tournaments</Link></li>
              <li><Link to=''>Engines</Link></li>
              <li><Link to=''>Books</Link></li>
              <li><Link to=''>Workers</Link></li>
            </ul>
            <ul className="nav navbar-nav navbar-right">
              <li><Link to=''><span className="glyphicon glyphicon-user"></span> Sign Up</Link></li>
              <li><Link to=''><span className="glyphicon glyphicon-log-in"></span> Login</Link></li>
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
        let term = pathTerms[i].charAt(0).toUpperCase() + pathTerms[i].slice(1);
        if (i < pathTerms.length - 1) {
          crumbItems.push(<li key={term}><a href="">{term}</a></li>);
        } else {
        crumbItems.push(<li className="active" key={term}>{term}</li>);
        }
      }
    }

    var link;
    if (this.props.navLink) {
      link = <li className="pull-right"><a href="">{this.props.navLink}<span className="glyphicon glyphicon-menu-right"></span></a></li>;
    }
    return (
      <ul className="breadcrumb">
        { crumbItems }
        { link }
      </ul>
    );
  }
}
