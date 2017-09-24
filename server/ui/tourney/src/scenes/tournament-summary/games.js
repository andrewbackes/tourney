import React, { Component } from 'react';

export default class GameTable extends Component {
  render() {
    let rows = [];
    this.props.runningGames.forEach( (game) => {
      rows.push(<GameTableRow key={game.id} game={game} history={this.props.history}/>)
    })
    return (
      <table className="table table-hover table-condensed">
        <thead>
          <tr>
            <th>Round</th>
            <th>White</th>
            <th>Black</th>
            <th>Result</th>
            <th>Winner</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          { rows }
        </tbody>
      </table>
    );
  }
}

function engineLabel(engine) {
  return engine.name + " " + engine.version;
}

class GameTableRow extends Component {
  handleClick(e) {
    this.props.history.push('/tournaments/' + this.props.game.tournamentId + '/games/' + this.props.game.id);
  }

  render() {
    return (
      <tr className='clickable' onClick={this.handleClick.bind(this)}>
        <td>{this.props.game.round}</td>
        <td>{engineLabel(this.props.game.contestants["0"])}</td> 
        <td>{engineLabel(this.props.game.contestants["1"])}</td> 
        <td>-</td>
        <td>-</td>
        <td>{this.props.game.status}</td>
      </tr>
    );
  }
}