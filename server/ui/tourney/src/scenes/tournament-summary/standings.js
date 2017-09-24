import React, { Component } from 'react';

export default class StandingsTable extends Component {
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
            <td></td>
            <td></td> 
            <td>-</td>
          </tr>
        </tbody>
      </table>
    );
  }
}