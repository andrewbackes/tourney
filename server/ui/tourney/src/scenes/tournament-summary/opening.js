import React, { Component } from 'react';

export default class OpeningBookTable extends Component {
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