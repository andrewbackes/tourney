import React, { Component } from 'react';
import Panel from 'components/panel';
import TournamentService from 'services/tournament';

class TournamentsDashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      active: TournamentService.getTournamentList(),
      pending: TournamentService.getTournamentList(),
      complete: TournamentService.getTournamentList()
    };
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

  refreshList() {
    this.setState({
      active: TournamentService.getTournamentList(),
      pending: TournamentService.getTournamentList(),
      complete: TournamentService.getTournamentList()
    });
  }

  render() {
    return (
      <div>
        <div className="row">
          <div className="col-xs-8">
            <Panel title="Active" mode="success" content={<ActiveTournamentsTable list={this.state.active}/>}/>
          </div>
          <div className="col-xs-4">
            <Panel title="Pending" mode="info" content={<PendingTournamentsTable list={this.state.pending}/>}/>
          </div>
        </div>
        <div className="row">
          <div className="col-xs-12">
            <Panel title="Completed" mode="default" content={<CompletedTournamentsTable list={this.state.complete}/>}/>
          </div>
        </div>
      </div>
    );
  }
}

export default TournamentsDashboard;

class ActiveTournamentsTable extends Component {
  render() {
    var rows = [];
    this.props.list.forEach( (tournament) => {
      rows.push(<ActiveTournamentsTableRow key={tournament.id} tournament={tournament}></ActiveTournamentsTableRow>);
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

class ActiveTournamentsTableRow extends Component {
  render() {
    return (
      <tr>
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
      rows.push(<PendingTournamentsTableRow key={tournament.id} tournament={tournament}></PendingTournamentsTableRow>);
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
  render() {
    return (
      <tr>
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
      rows.push(<CompletedTournamentsTableRow key={tournament.id} tournament={tournament}></CompletedTournamentsTableRow>);
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
  render() {
    return (
      <tr>
        <td>{this.props.tournament.id}</td>
        <td>{formatTimeControl(this.props.tournament.settings.timeControl)}</td>
        <td>-</td>
      </tr>
    )
  }
}