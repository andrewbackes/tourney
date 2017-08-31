import React, { Component } from 'react';

class GameList extends Component {
  render() {
    return (
      <div className="panel panel-default">
        <div className="panel-body">
          <GameTable/>
        </div>
      </div>
    );
  }
}

class GameTable extends Component {
  render() {
    return (
      <table className="table table-hover table-condensed">
        <thead>
          <tr>
            <th>Round</th>
            <th>White</th>
            <th>Black</th>
            <th>Result</th>
            <th>Winner</th>
            <th>Ending Condistion</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>1</td>
            <td>Dirty-Bit 09ba34ef</td> 
            <td>Dirty-Bit 1ab34ef</td>
            <td>1-0</td>
            <td>Dirty-Bit 09ba34ef</td>
            <td>Checkmate</td>
          </tr>
          <tr>
            <td>2</td>
            <td>Dirty-Bit 1ab34ef</td> 
            <td>Dirty-Bit 09ba34ef</td> 
            <td>0-1</td>
            <td>Dirty-Bit 09ba34ef</td>
            <td>Checkmate</td>
          </tr>
        </tbody>
      </table>
    );
  }
}

export default GameList;