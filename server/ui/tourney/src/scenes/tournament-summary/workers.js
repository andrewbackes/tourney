import React, { Component } from 'react';

export default class WorkersTable extends Component {
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