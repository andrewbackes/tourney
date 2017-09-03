import React, { Component } from 'react';
import { Switch, Route, Redirect } from 'react-router-dom'

import TournamentsDashboard from 'components/tournaments';
import TournamentDashboard from 'components/tournament';
import GameList from 'components/games';
import Game from 'components/game';

import WorkersDashboard from 'components/workers';
import EnginesDashboard from 'components/engines';
import BooksDashboard from 'components/books';

import SignupScreen from 'components/signup';
import LoginScreen from 'components/login';

import NotFoundScreen from 'components/notfound';

export default class Main extends Component {
  render() {
    return (
      <main>
        <Switch>
          <Redirect exact from="/" to="/tournaments"/>
          <Route exact path='/tournaments/:tournamentId/games/:gameId' component={Game}/>
          <Route exact path='/tournaments/:tournamentId/games' component={GameList}/>
          <Route exact path='/tournaments/:tournamentId' component={TournamentDashboard}/>
          <Route exact path='/tournaments' component={TournamentsDashboard}/>
          <Route exact path='/workers' component={WorkersDashboard}/>
          <Route exact path='/books' component={BooksDashboard}/>
          <Route exact path='/engines' component={EnginesDashboard}/>
          <Route exact path='/signup' component={SignupScreen}/>
          <Route exact path='/login' component={LoginScreen}/>
          <Route component={NotFoundScreen} />
        </Switch>
      </main>
    );
  }
}