import React, { Component } from 'react';
import Panel from 'components/panel';
import TournamentService from 'services/tournament';

import RunningTournamentsTable from 'scenes/tournament-list/running';
import PendingTournamentsTable from 'scenes/tournament-list/pending';
import CompletedTournamentsTable from 'scenes/tournament-list/completed';

import 'style/main.css';

export default class TournamentsDashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      running: [],
      pending: [],
      complete: []
    };
    this.setRunningTournaments = this.setRunningTournaments.bind(this);
    this.setPendingTournaments = this.setPendingTournaments.bind(this);
    this.setCompleteTournaments = this.setCompleteTournaments.bind(this);
    this.refreshList();
  }

  componentDidMount() {
    this.timerID = setInterval(
      () => this.refreshList(),
      1000
    );
  }

  componentWillUnmount() {
    clearInterval(this.timerID);
  }

  setRunningTournaments(tournaments) {
    if (this.timerID) {
      this.setState({ running: tournaments });
    }
  }

  setPendingTournaments(tournaments) {
    if (this.timerID) {
      this.setState({ pending: tournaments });
    }
  }

  setCompleteTournaments(tournaments) {
    if (this.timerID) {
      this.setState({ complete: tournaments });
    }
  }

  refreshList() {
    TournamentService.getTournamentList('running', this.setRunningTournaments);
    TournamentService.getTournamentList('pending', this.setPendingTournaments);
    TournamentService.getTournamentList('complete', this.setCompleteTournaments);
  }

  render() {
    return (
      <div>
        <div className="row">
          <div className="col-xs-8">
            <Panel title="Running" mode="success" content={<RunningTournamentsTable list={this.state.running} history={this.props.history}/>}/>
          </div>
          <div className="col-xs-4">
            <Panel title="Pending" mode="info" content={<PendingTournamentsTable list={this.state.pending} history={this.props.history}/>}/>
          </div>
        </div>
        <div className="row">
          <div className="col-xs-12">
            <Panel title="Completed" mode="default" content={<CompletedTournamentsTable list={this.state.complete} history={this.props.history}/>}/>
          </div>
        </div>
      </div>
    );
  }
}
