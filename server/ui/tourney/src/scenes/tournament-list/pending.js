import React, { Component } from 'react';

import TimeControl from 'util/time-control';

export default class PendingTournamentsTable extends Component {
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
        <td>{TimeControl.format(this.props.tournament.settings.timeControl)}</td>
      </tr>
    )
  }
}