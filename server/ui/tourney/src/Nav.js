import React, { Component } from 'react';

class NavBar extends Component {
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
              <li className={this.props.navPath[0].toLowerCase() === 'tournaments' ? 'active' : ''}><a href="">Tournaments</a></li>
              <li className={this.props.navPath[0].toLowerCase() === 'engines' ? 'active' : ''}><a href="">Engines</a></li>
              <li className={this.props.navPath[0].toLowerCase() === 'books' ? 'active' : ''}><a href="">Books</a></li>
              <li className={this.props.navPath[0].toLowerCase() === 'workers' ? 'active' : ''}><a href="">Workers</a></li>
            </ul>
            <ul className="nav navbar-nav navbar-right">
              <li><a href=""><span className="glyphicon glyphicon-user"></span> Sign Up</a></li>
              <li><a href=""><span className="glyphicon glyphicon-log-in"></span> Login</a></li>
            </ul>
            </div>
          </div>
        </nav>
        <NavBreadcrumbs navPath={this.props.navPath} navLink={this.props.navLink}/>
      </div>
    ); 
  }
}

export default NavBar;

class NavBreadcrumbs extends Component {
  render() {
    var crumbItems = [];
    this.props.navPath.forEach( (item) => {
      crumbItems.push(<li key={item}><a href="">{item}</a></li>);
    });
    var lastIndex = this.props.navPath.length - 1;
    crumbItems[lastIndex] = <li className="active" key={this.props.navPath[lastIndex]}>{this.props.navPath[lastIndex]}</li>;
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