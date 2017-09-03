import React, { Component } from 'react';
import Panel from 'components/panel';
import TournamentService from 'services/tournament';
import 'style/main.css';

class TournamentsDashboard extends Component {
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
      10000
    );
  }

  componentWillUnmount() {
    clearInterval(this.timerID);
  }

  setRunningTournaments(tournaments) {
    this.setState({ running: tournaments });
  }

  setPendingTournaments(tournaments) {
    this.setState({ pending: tournaments });
  }

  setCompleteTournaments(tournaments) {
    this.setState({ complete: tournaments });
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

export default TournamentsDashboard;

class RunningTournamentsTable extends Component {
  render() {
    var rows = [];
    this.props.list.forEach( (tournament) => {
      rows.push(<RunningTournamentsTableRow key={tournament.id} tournament={tournament} history={this.props.history}></RunningTournamentsTableRow>);
    });
    return (
      <table className="table table-hover table-condensed">
        <thead>
          <tr>
            <th>Name</th>
            <th>Time Control</th>
            <th>Leader</th>
            <th>Progress</th>
          </tr>
        </thead>
        <tbody>
        { rows }
        </tbody>
      </table>
    );
  }
}

function formatTimeControl(timeControl) {
  let format = (val) => {
    if (val) {
      return val.toString();
    } else {
      return "";
    }
  };
  return format(timeControl.moves) + "/" + format(timeControl.time) + "+" + format(timeControl.increment);
}

class RunningTournamentsTableRow extends Component {
  handleClick(e) {
    this.props.history.push('/tournaments/' + this.props.tournament.id);
  }

  render() {
    return (
      <tr className='clickable' onClick={this.handleClick.bind(this)}>
        <td>{this.props.tournament.id}</td>
        <td>{formatTimeControl(this.props.tournament.settings.timeControl)}</td>
        <td>-</td>
        <td>-</td>
      </tr>
    )
  }
}
  
class PendingTournamentsTable extends Component {
  render() {
    var rows = [];
    this.props.list.forEach( (tournament) => {
      rows.push(<PendingTournamentsTableRow key={tournament.id} tournament={tournament} history={this.props.history}></PendingTournamentsTableRow>);
    });
    return (
      <table className="table table-hover table-condensed">
        <thead>
          <tr>
            <th>Name</th>
            <th>Time Control</th>
          </tr>
        </thead>
        <tbody>
          { rows }
        </tbody>
      </table>
    );
  }
}

class PendingTournamentsTableRow extends Component {
  handleClick(e) {
    this.props.history.push('/tournaments/' + this.props.tournament.id);
  }

  render() {
    return (
      <tr className='clickable' onClick={this.handleClick.bind(this)}>
        <td>{this.props.tournament.id}</td>
        <td>{formatTimeControl(this.props.tournament.settings.timeControl)}</td>
      </tr>
    )
  }
}

class CompletedTournamentsTable extends Component {
  render() {
    var rows = [];
    this.props.list.forEach( (tournament) => {
      rows.push(<CompletedTournamentsTableRow key={tournament.id} tournament={tournament} history={this.props.history}></CompletedTournamentsTableRow>);
    });
    return (
      <table className="table table-hover table-condensed">
        <thead>
          <tr>
            <th>Name</th>
            <th>Time Control</th>
            <th>Winner</th>
          </tr>
        </thead>
        <tbody>
          { rows }
        </tbody>
      </table>
    );
  }
}

class CompletedTournamentsTableRow extends Component {
  handleClick(e) {
    this.props.history.push('/tournaments/' + this.props.tournament.id);
  }

  render() {
    return (
      <tr className='clickable' onClick={this.handleClick.bind(this)}>
        <td>{this.props.tournament.id}</td>
        <td>{formatTimeControl(this.props.tournament.settings.timeControl)}</td>
        <td>-</td>
      </tr>
    )
  }
}