import React, { Component } from 'react';
import { Link } from 'react-router-dom'

export default class NavBreadcrumbs extends Component {
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