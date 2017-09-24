import React, { Component } from 'react';
import { Link } from 'react-router-dom'
import Panel from 'components/panel';
import TournamentService from 'services/tournament';

import OpeningBookTable from 'scenes/tournament-summary/opening';
import StandingsTable from 'scenes/tournament-summary/standings';
import GameTable from 'scenes/tournament-summary/games';
import MathupsTable from 'scenes/tournament-summary/matchups';
import TimeControlTable from 'scenes/tournament-summary/time-control';
import WorkersTable from 'scenes/tournament-summary/workers';

export default class TournamentDashboard extends Component {
  constructor(props) {
    super(props);
    this.state = { 
      tournament: {},
      runningGames: [],
      workers: []
    };
    this.setTournament = this.setTournament.bind(this);
    this.setRunningGames = this.setRunningGames.bind(this);
    this.refreshTournament();
  }

  componentDidMount() {
    if (this.state.tournament.status !== "Complete") {
      this.timerID = setInterval(
        () => this.refreshTournament(),
        1000
      );
    }
  }

  componentWillUnmount() {
    clearInterval(this.timerID);
  }

  setTournament(tournament) {
    if (this.timerID) {
      this.setState({ tournament: tournament });
    }
    if (this.state.tournament.status === "Complete") {
      clearInterval(this.timerID);
    }
  }

  setRunningGames(games) {
    if (this.timerID) {
      this.setState({ runningGames: games });
    }
  }

  refreshTournament() {
    TournamentService.getTournament(this.props.match.params.tournamentId, this.setTournament)
    TournamentService.getGameList(this.props.match.params.tournamentId, this.setRunningGames, "running")
  }

  render() {
    return (
      <div>
        <div className="row">
          <div className="col-xs-12">
            <Panel title="Running Games" mode="success" content={
              <div>
                <GameTable runningGames={this.state.runningGames} history={this.props.history}/>
                <div className="panel-body text-right">
                  <Link to={'/tournaments/' + this.props.match.params.tournamentId + '/games'}>All Games<span className="glyphicon glyphicon-menu-right"></span></Link>
                </div>
              </div>
            }/>
          </div>
        </div>
        <div className="row">
          <div className="col-xs-6">
            <Panel title="Standings" mode="default" content={<StandingsTable tournament={this.props.tournament}/>}/>
          </div>
          <div className="col-xs-6">
            <Panel title="Matchups" mode="default" content={<MathupsTable tournament={this.props.tournament}/>}/>
          </div>
        </div>
        { this.state.tournament && this.state.tournament.settings &&
        <div className="row">
          <div className="col-xs-4">
            <Panel title="Time Control" mode="default" content={<TimeControlTable timeControl={this.state.tournament.settings.timeControl}/>}/>
          </div>
          <div className="col-xs-4">
            <Panel title="Opening Book" mode="default" content={<OpeningBookTable opening={this.state.tournament.settings.opening}/>}/>
          </div>
          <div className="col-xs-4">
            <Panel title="Workers" mode="success" content={<WorkersTable workers={this.state.workers}/>}/>
          </div>
        </div>
        }
      </div>
    );
  }
}
