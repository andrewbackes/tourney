import React, { Component } from 'react';
import Panel from 'components/panel';
import TournamentService from 'services/tournament';
import Board from 'components/chessboard';

export default class GameDashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      game: {}
    };
    this.setGame = this.setGame.bind(this);
    this.refreshGame()
  }

  componentDidMount() {
    if (this.state.game && this.state.game.status !== "Complete") {
      this.timerID = setInterval(
        () => this.refreshGame(),
        10000
      );
    }
  }

  componentWillUnmount() {
    clearInterval(this.timerID);
  }

  setGame(game) {
    this.setState({ game: game });
    if (this.state.game.status === "Complete") {
      clearInterval(this.timerID);
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
            <Panel title="Board" mode="default" content={<Board fen="rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"/>}/>
          </div>
          <div className="col-xs-4">
            <Panel title="Moves" mode="default" content={<MoveTable game={this.state.game}/>}/>
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
        rows.push(<MoveTableRow key={pos.fen} lastMove={pos.lastMove}/>)
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
  render() {
    return (
      <tr>
        <td></td>
        <td>{this.props.lastMove.source}->{this.props.lastMove.destination}</td>
        <td>{this.props.lastMove.duration}</td>
      </tr>
    );
  }
}
