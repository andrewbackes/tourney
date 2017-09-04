import React, { Component } from 'react';
import Panel from 'components/panel';
import TournamentService from 'services/tournament';
import Board from 'components/chessboard';

export default class GameDashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      game: {},
      fen: ""
    };
    this.initGame = this.initGame.bind(this);
    this.setGame = this.setGame.bind(this);
    this.setFen = this.setFen.bind(this);
    TournamentService.getGame(this.props.match.params.tournamentId, this.props.match.params.gameId, this.initGame);
  }

  componentDidMount() {
    this.timerID = setInterval(
      () => this.refreshGame(),
      500
    );
  }

  componentWillUnmount() {
    clearInterval(this.timerID);
  }

  initGame(game) {
    this.setState({ game: game });
    if (game.status === "Complete") {
      clearInterval(this.timerID);
      this.setFen(game.positions[0].fen);
    } else {
      this.setFen(this.state.game.positions[this.state.game.positions.length-1].fen);
    }
  }

  setGame(game) {
    let updateFen = false;
    if (this.state.game && this.state.game.positions) {
      updateFen = this.state.fen === this.state.game.positions[this.state.game.positions.length-1].fen
    }
    if (this.timerID) {
      this.setState({ game: game });
    }
    if (updateFen) {
      this.setFen(this.state.game.positions[this.state.game.positions.length-1].fen);
    }
    
    if (this.state.game.status === "Complete") {
      clearInterval(this.timerID);
    }
  }

  setFen(fen) {
    if (this.timerID) {
      this.setState({ fen: fen });
    }
  }

  refreshGame() {
    TournamentService.getGame(this.props.match.params.tournamentId, this.props.match.params.gameId, this.setGame)
  }
  
  render() {
    return (
      <div>
        <div className="row">
          <div className="col-xs-8">
            <Panel title="Board" mode="default" content={<Board fen={this.state.fen}/>}/>
          </div>
          <div className="col-xs-4">
            <Panel title="Moves" mode="default" content={ <MoveTable game={this.state.game} setFen={this.setFen} currentFen={this.state.fen} /> }/>
          </div>
        </div>
      </div>
    );
  }
}

class MoveTable extends Component {
  render() {
    let rows = [];
    if (this.props.game.positions) {
      this.props.game.positions.forEach( (pos) => {
        rows.push(<MoveTableRow key={pos.fen} fen={pos.fen} lastMove={pos.lastMove} setFen={this.props.setFen} currentFen={this.props.currentFen}/>)
      });
    }
    return (
      <table className="table table-hover table-condensed">
        <thead>
          <tr>
            <th>Number</th>
            <th>Move</th>
            <th>Duration</th>
          </tr>
        </thead>
        <tbody>
          { rows }
        </tbody>
      </table>
    );
  }
}

class MoveTableRow extends Component {
  handleClick(e) {
    this.props.setFen(this.props.fen);
  }
  render() {
    let active = this.props.fen === this.props.currentFen;
    return (
      <tr className={'clickable ' + (active ? 'active' : '')} onClick={this.handleClick.bind(this)}>
        <td></td>
        <td>{this.props.lastMove.source}->{this.props.lastMove.destination}</td>
        <td>{this.props.lastMove.duration}</td>
      </tr>
    );
  }
}
