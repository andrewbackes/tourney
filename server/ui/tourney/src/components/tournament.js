import React, { Component } from 'react';
import Panel from 'components/panel';
import TournamentService from 'services/tournament';

export default class TournamentDashboard extends Component {
  constructor(props) {
    super(props);
    this.state = { 
      tournament: TournamentService.getTournament("something"),
      workers: []
    };
  }

  componentDidMount() {
    if (this.state.tournament.status !== "Complete") {
      this.timerID = setInterval(
        () => this.refreshTournament(),
        10000
      );
    }
  }

  componentWillUnmount() {
    clearInterval(this.timerID);
  }

  refreshTournament() {
    this.setState({
      tournament: TournamentService.getTournament("something")
    });
    if (this.state.tournament.status === "Complete") {
      clearInterval(this.timerID);
    }
  }

  render() {
    return (
      <div>
        <div className="row">
          <div className="col-xs-6">
            <Panel title="Standings" mode="default" content={<StandingsTable tournament={this.props.tournament}/>}/>
          </div>
          <div className="col-xs-6">
            <Panel title="Matchups" mode="default" content={<MathupsTable tournament={this.props.tournament}/>}/>
          </div>
        </div>
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
      </div>
    );
  }
}

class StandingsTable extends Component {
  render() {
    return (
      <table className="table table-condensed">
        <thead>
          <tr>
            <th>Position</th>
            <th>Name</th>
            <th>Score</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>1</td>
            <td>Dirty-Bit 09ba34ef</td> 
            <td>-</td>
          </tr>
          <tr>
            <td>2</td>
            <td>Dirty-Bit 1ab34ef</td> 
            <td>-</td>
          </tr>
        </tbody>
      </table>
    );
  }
}

class MathupsTable extends Component {
  render() {
    return (
      <table className="table table-condensed">
        <thead>
          <tr>
            <th>Engine</th>
            <th>Opponent</th>
            <th>Score</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Dirty-Bit 09ba34ef</td> 
            <td>Dirty-Bit 1ab34ef</td>
            <td>6-3-1</td>
          </tr>
          <tr>
            <td>Dirty-Bit 1ab34ef</td>
            <td>Dirty-Bit 09ba34ef</td>
            <td>3-6-1</td>
          </tr>
        </tbody>
      </table>
    );
  }
}

class TimeControlTable extends Component {
  render() {
    return (
      <table className="table table-condensed">
        <tbody>
          <tr>
            <th>Time</th>
            <td>{ this.props.timeControl.time }</td>
          </tr>
          <tr>
            <th>Moves</th>
            <td>{ this.props.timeControl.moves }</td>
          </tr>
          <tr>
            <th>Bonus</th>
            <td>{ this.props.timeControl.increment }</td>
          </tr>
          <tr>
            <th>Repeating</th>
            <td>{ this.props.timeControl.repeating }</td>
          </tr>
        </tbody>
      </table>
    );
  }
}

class OpeningBookTable extends Component {
  render() {
    return (
      <table className="table table-condensed">
        <tbody>
          <tr><th>Name</th><td>{this.props.opening && this.props.opening.bookName ? this.props.opening.bookName : "-"}</td></tr>
          <tr><th>Depth</th><td>{this.props.opening && this.props.opening.depth ? this.props.opening.depth : "-"}</td></tr>
          <tr><th>Mirrored</th><td>{this.props.opening && this.props.opening.mirrored ? this.props.opening.mirrored : "-"}</td></tr>
          <tr><th>Randomized</th><td>{this.props.opening && this.props.opening.randomized ? this.props.opening.randomized : "-"}</td></tr>
        </tbody>
      </table>
    );
  }
}

class WorkersTable extends Component {
  render() {
    let rows = [];
    this.props.workers.forEach( (worker) => {
      rows.push(<WorkersTableRow key={worker.id} worker={worker}/>);
    });
    return (
      <table className="table table-hover table-condensed">
        <thead>
          <tr>
            <th>User</th>
            <th>Id</th>
            <th>Game</th>
          </tr>
        </thead>
        <tbody>
          {rows}
        </tbody>
      </table>
    );
  }
}

class WorkersTableRow extends Component {
  render() {
    return (
      <tr>
        <td>{this.props.worker.id}</td>
        <td>{this.props.worker.gameId}</td> 
      </tr>
    );
  }
}
