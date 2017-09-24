import React, { Component } from 'react';

import NavBar from 'layout/header/nav';
import NavBreadcrumbs from 'layout/header/breadcrumbs';

export default class Header extends Component {
  render() {
    return (
      <div>
        <NavBar location={this.props.location}/>
        <NavBreadcrumbs navPath={this.props.location.pathname}/>
      </div>
    ); 
  }
}
