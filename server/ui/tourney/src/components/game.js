import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Panel from 'components/panel';
import TournamentService from 'services/tournament';
import Board from 'components/chessboard';

export default class GameDashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      game: {},
      position: {
        fen: ""
      }
    };
    this.initGame = this.initGame.bind(this);
    this.setGame = this.setGame.bind(this);
    this.setPosition = this.setPosition.bind(this);
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
      this.setPosition(game.positions[0]);
    } else {
      this.setPosition(this.state.game.positions[this.state.game.positions.length-1]);
    }
  }

  setGame(game) {
    let updatePosition = false;
    if (this.state.game && this.state.game.positions) {
      updatePosition = this.state.position === this.state.game.positions[this.state.game.positions.length-1]
    }
    if (this.timerID) {
      this.setState({ game: game });
    }
    if (updatePosition) {
      this.setPosition(this.state.game.positions[this.state.game.positions.length-1]);
    }
    
    if (this.state.game.status === "Complete") {
      clearInterval(this.timerID);
    }
  }

  setPosition(position) {
    if (this.timerID) {
      this.setState({ position: position });
    }
  }

  refreshGame() {
    TournamentService.getGame(this.props.match.params.tournamentId, this.props.match.params.gameId, this.setGame)
  }
  
  render() {
    let mode = 'default';
    if (this.state.game.status) {
      if (this.state.game.status.toLowerCase() === "running") {
        mode = 'success';
      } else if (this.state.game.status.toLowerCase() === "pending") {
        mode = 'info';
      }
    }
    return (
      <div>
        <div className="row">
          <div className="col-xs-8">
            <Panel title="Board" mode={ mode } content={<Board position={this.state.position}/>}/>
          </div>
          <div className="col-xs-4">
            <Panel title="Moves" mode="default" content={ <MoveTable game={this.state.game} setPosition={this.setPosition} currentPosition={this.state.position} /> }/>
          </div>
        </div>
      </div>
    );
  }
}

class MoveTable extends Component {
  

  scrollToBottom = () => {
    const node = ReactDOM.findDOMNode(this.tbody);
    if (node !== null && node.lastElementChild !== null) {
      if (node.children.length > 24) {
        if (node.lastElementChild.classList && node.lastElementChild.classList.contains('active')) {
          if (node.lastElementChild.lastElementChild) {
            node.lastElementChild.lastElementChild.scrollIntoView({ behavior: "smooth" });
          }
        }
      }
    }
  }
  
  componentDidMount() {
    this.scrollToBottom();
  }
  
  componentDidUpdate() {
    this.scrollToBottom();
  }


  render() {
    let rows = [];
    if (this.props.game.positions) {
      this.props.game.positions.forEach( (pos) => {
        rows.push(<MoveTableRow key={pos.fen} position={pos} lastMove={pos.lastMove} setPosition={this.props.setPosition} currentPosition={this.props.currentPosition}/>)
      });
    }
    return (
      <div>
        <table className="table table-hover table-condensed table-fixed">
          <thead>
            <tr>
              <th className="col-xs-4">Number</th>
              <th className="col-xs-4">Move</th>
              <th className="col-xs-4">Duration</th>
            </tr>
          </thead>
          <tbody ref={(el) => { this.tbody = el; }}>
            { rows }
          </tbody>
        </table>
      </div>
    );
  }
}

class MoveTableRow extends Component {
  handleClick(e) {
    this.props.setPosition(this.props.position);
  }
  render() {
    let active = this.props.position.fen === this.props.currentPosition.fen;
    return (
      <tr className={'clickable ' + (active ? 'active' : '')} onClick={this.handleClick.bind(this)}>
        <td className="col-xs-4">-</td>
        <td className="col-xs-4">{this.props.lastMove.source}->{this.props.lastMove.destination}</td>
        <td className="col-xs-4">{this.props.lastMove.duration ? this.props.lastMove.duration : "-"}</td>
      </tr>
    );
  }
}
