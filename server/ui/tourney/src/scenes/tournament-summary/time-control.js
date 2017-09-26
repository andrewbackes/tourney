import React, { Component } from 'react';

import Duration from 'util/duration';

export default class TimeControlTable extends Component {
  render() {
    return (
      <table className="table table-condensed">
        <tbody>
          <tr>
            <th>Time</th>
            <td>{ this.props.timeControl.time ? Duration.format(this.props.timeControl.time) : "-" }</td>
          </tr>
          <tr>
            <th>Moves</th>
            <td>{ this.props.timeControl.moves ? this.props.timeControl.moves : "-" }</td>
          </tr>
          <tr>
            <th>Bonus</th>
            <td>{ this.props.timeControl.increment ? Duration.format(this.props.timeControl.increment) : "-" }</td>
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