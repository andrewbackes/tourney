import React, { Component } from 'react';
import Panel from 'components/panel';
import TournamentService from 'services/tournament';

export default class GameList extends Component {
  constructor(props) {
    super(props);
    this.state = {
      tournament: {},
      gameList: [],
      filterText: ""
    };
    this.handleFilterTextInput = this.handleFilterTextInput.bind(this);
    this.setTournament = this.setTournament.bind(this);
    this.setGameList = this.setGameList.bind(this);
    this.refreshState()
  }

  componentDidMount() {
    if (this.state.tournament.status !== "Complete") {
      this.timerID = setInterval(
        () => this.refreshState(),
        10000
      );
    }
  }

  componentWillUnmount() {
    clearInterval(this.timerID);
  }

  handleFilterTextInput(filterText) {
    this.setState({
      filterText: filterText
    });
  }
  
  setTournament(tournament) {
    this.setState({ tournament: tournament });
    if (this.state.tournament.status === "Complete") {
      clearInterval(this.timerID);
    }
  }

  setGameList(gameList) {
    this.setState({ gameList: gameList });
  }

  refreshState() {
    TournamentService.getTournament(this.props.match.params.tournamentId, this.setTournament)
    TournamentService.getGameList(this.props.match.params.tournamentId, this.setGameList)
  }
  
  render() {
    return (
      <div>
        <Panel title="Search" mode="default" content={
          <Search onFilterTextInput={this.handleFilterTextInput} filterText={this.state.filterText}/>
        }/>
        <div className="panel panel-default">
          <div className="panel-body">
            <GameTable gameList={this.state.gameList} filterText={this.state.filterText} history={this.props.history}/>
          </div>
        </div>
      </div>
    );
  }
}

class GameTable extends Component {
  
  filterGame(game, filterText) {
    if (filterText === "") {
      return true;
    }
    return JSON.stringify(game).toLowerCase().includes(filterText);
  }

  render() {
    let rows = [];
    this.props.gameList.forEach( (game) => {
      if (this.filterGame(game, this.props.filterText)) {
        rows.push(<GameTableRow key={game.id} game={game} history={this.props.history}/>)
      }
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
class Search extends Component {
  constructor(props) {
    super(props);
    this.handleFilterTextInputChange = this.handleFilterTextInputChange.bind(this);
  }

  handleFilterTextInputChange(e) {
    this.props.onFilterTextInput(e.target.value);
  }

  render() {
    return (
    <div className="input-group">
      <span className="input-group-addon">Filter</span>
      <input
            type="text"
            className="form-control"
            placeholder="Search..."
            value={this.props.filterText}
            onChange={this.handleFilterTextInputChange}
        />
    </div>
    );
  }
}