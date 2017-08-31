import React, { Component } from 'react';
import './App.css';
import NavBar from './Nav'
import TournamentsDashboard from './Tournaments'

class App extends Component {
  content(navPath) {
    if (navPath[0].toLowerCase() === 'tournaments') {
      if (navPath.length === 1) {
        return <TournamentsDashboard/>;
      }
    }
  }

  render() {
    var navPath = ["Tournaments", "123"];
    return (
      <div className="App">
        <div className="col-xs-10 col-xs-offset-1">
          <NavBar navPath={navPath}/>
          {this.content(navPath)}
        </div>
      </div>
    );
  }
}

export default App;
