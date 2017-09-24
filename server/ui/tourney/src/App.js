import React, { Component } from 'react';
import {withRouter} from 'react-router-dom';
import Header from 'layout/header';
import Main from 'layout/main';
import 'style/main.css';

class App extends Component {

  render() {
    const RoutedHeader = withRouter(props => <Header {...props}/>);
    return (
      <div className="col-xs-10 col-xs-offset-1">
        <RoutedHeader/>
        <Main/>
      </div>
    );
  }
}

export default App;
