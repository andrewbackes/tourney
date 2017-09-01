import React, { Component } from 'react';
import 'style/App.css';
import NavBar from 'components/nav';
import TournamentsDashboard from 'components/tournaments';
import TournamentDashboard from 'components/tournament';
import GameList from 'components/games';

class App extends Component {
  content(navPath) {
    if (navPath[0].toLowerCase() === 'tournaments') {
      if (navPath.length === 1) {
        return <TournamentsDashboard/>;
      } else if (navPath.length === 2) {
        return <TournamentDashboard/>;
      } else if (navPath.length === 3) {
        return <GameList/>;
      }
    }
  }

  navLink(navPath) {
    if (navPath[0].toLowerCase() === 'tournaments' && navPath.length === 2) {
      return 'Games';
    }
  }
  
  render() {
    var navPath = ["Tournaments", "123"];
    return (
      <div className="App">
        <div className="col-xs-10 col-xs-offset-1">
          <NavBar navPath={navPath} navLink={this.navLink(navPath)}/>
          {this.content(navPath)}
        </div>
      </div>
    );
  }
}

export default App;
