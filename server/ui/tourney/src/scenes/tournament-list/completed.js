import React, { Component } from 'react';

import TimeControl from 'util/time-control';

export default class CompletedTournamentsTable extends Component {
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
        <td>{TimeControl.format(this.props.tournament.settings.timeControl)}</td>
        <td>-</td>
      </tr>
    )
  }
}