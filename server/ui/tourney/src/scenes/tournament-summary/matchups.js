import React, { Component } from 'react';

export default class MathupsTable extends Component {
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
            <td></td> 
            <td></td>
            <td></td>
          </tr>
        </tbody>
      </table>
    );
  }
}