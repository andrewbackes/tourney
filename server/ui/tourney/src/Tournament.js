import React, { Component } from 'react';
import Panel from './Panel'

class TournamentDashboard extends Component {
  render() {
    return (
      <div>
        <div className="row">
          <div className="col-xs-6">
            <Panel title="Standings" mode="default" content={<StandingsTable/>}/>
          </div>
          <div className="col-xs-6">
            <Panel title="Matchups" mode="default" content={<MathupsTable/>}/>
          </div>
        </div>
        <div className="row">
          <div className="col-xs-4">
            <Panel title="Time Control" mode="default" content={<TimeControlTable/>}/>
          </div>
          <div className="col-xs-4">
            <Panel title="Opening Book" mode="default" content={<OpeningBookTable/>}/>
          </div>
          <div className="col-xs-4">
            <Panel title="Workers" mode="success" content={<WorkersTable/>}/>
          </div>
        </div>
      </div>
    );
  }
}

export default TournamentDashboard;

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
            <td>5min</td>
          </tr>
          <tr>
            <th>Moves</th>
            <td>40</td>
          </tr>
          <tr>
            <th>Bonus</th>
            <td>10s</td>
          </tr>
          <tr>
            <th>Repeating</th>
            <td>True</td>
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
        <tr><th>Name</th><td>2800elo</td></tr>
        <tr><th>Depth</th><td>4</td></tr>
        <tr><th>Mirrored</th><td>True</td></tr>
        <tr><th>Randomized</th><td>True</td></tr>
        </tbody>
      </table>
    );
  }
}

class WorkersTable extends Component {
  render() {
    return (
      <table className="table table-hover table-condensed">
        <thead>
          <tr>
            <th>User</th>
            <th>Game</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>andrewbackes</td>
            <td>Round 1</td> 
          </tr>
          <tr>
            <td>andrewbackes</td>
            <td>Round 2</td> 
          </tr>
        </tbody>
      </table>
    );
  }
}
