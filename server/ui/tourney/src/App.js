import React, { Component } from 'react';
import {withRouter} from 'react-router-dom';
import NavBar from 'components/nav';
import Main from 'components/main';

class App extends Component {

  render() {
    const RoutedNavBar = withRouter(props => <NavBar {...props}/>);
    return (
      <div className="col-xs-10 col-xs-offset-1">
        <RoutedNavBar/>
        <Main/>
      </div>
    );
  }
}

export default App;
