import React, { Component } from 'react';
import Panel from 'components/panel';

class TournamentsDashboard extends Component {
  render() {
    return (
      <div>
        <div className="row">
          <div className="col-xs-8">
            <Panel title="Active" mode="success" content={<ActiveTournamentsTable/>}/>
          </div>
          <div className="col-xs-4">
            <Panel title="Pending" mode="info" content={<PendingTournamentsTable/>}/>
          </div>
        </div>
        <div className="row">
          <div className="col-xs-12">
            <Panel title="Completed" mode="default" content={<CompletedTournamentsTable/>}/>
          </div>
        </div>
      </div>
    );
  }
}

export default TournamentsDashboard;

class ActiveTournamentsTable extends Component {
  render() {
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
          <tr>
            <td>Dirty-Bit 1ab34ef Evaluation</td>
            <td>5/40+10</td>
            <td>Dirty-Bit 09ba34ef</td>
            <td>-</td>
          </tr>
          <tr>
            <td>Dirty-Bit 09ba34ef Evaluation</td>
            <td>5/40+10</td>
            <td>Dirty-Bit 09ba34ef</td>
            <td>-</td>
          </tr>
        </tbody>
      </table>
    );
  }
}
  
class PendingTournamentsTable extends Component {
  render() {
    return (
      <table className="table table-hover table-condensed">
        <thead>
          <tr>
            <th>Name</th>
            <th>Time Control</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Dirty-Bit 97ba31f Evaluation</td>
            <td>5/40+10</td>
          </tr>
          <tr>
            <td>Dirty-Bit 19ba34ef Evaluation</td>
            <td>5/40+10</td>
          </tr>
        </tbody>
      </table>
    );
  }
}

class CompletedTournamentsTable extends Component {
  render() {
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
          <tr>
            <td>Dirty-Bit 1ab34ef Evaluation</td>
            <td>5/40+10</td>
            <td>Dirty-Bit 09ba34ef</td>
          </tr>
          <tr>
            <td>Dirty-Bit 09ba34ef Evaluation</td>
            <td>5/40+10</td>
            <td>Dirty-Bit 09ba34ef</td>
          </tr>
          <tr>
            <td>Dirty-Bit 97ba31f Evaluation</td>
            <td>5/40+10</td>
            <td>Dirty-Bit 97ba31f</td>
          </tr>
        </tbody>
      </table>
    );
  }
}