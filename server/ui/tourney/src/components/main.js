import React, { Component } from 'react';
import { Switch, Route } from 'react-router-dom'

import TournamentsDashboard from 'components/tournaments';
import TournamentDashboard from 'components/tournament';
import GameList from 'components/games';
import Game from 'components/game';

export default class Main extends Component {
  render() {
    return (
      <main>
        <Switch>
          <Route exact path='/' component={TournamentsDashboard}/>
          <Route exact path='/tournaments/:tournamentId/games/:gameId' component={Game}/>
          <Route exact path='/tournaments/:tournamentId/games' component={GameList}/>
          <Route exact path='/tournaments/:tournamentId' component={TournamentDashboard}/>
          <Route exact path='/tournaments' component={TournamentsDashboard}/>
        </Switch>
      </main>
    );
  }
}