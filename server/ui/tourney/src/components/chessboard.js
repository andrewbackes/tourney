import React, { Component } from 'react';
import 'pieces/piece01.png';
import 'pieces/piece02.png';
import 'pieces/piece03.png';
import 'pieces/piece04.png';
import 'pieces/piece05.png';
import 'pieces/piece06.png';
import 'pieces/piece11.png';
import 'pieces/piece12.png';
import 'pieces/piece13.png';
import 'pieces/piece14.png';
import 'pieces/piece15.png';
import 'pieces/piece16.png';

export default class Board extends Component {
  render() {
    let squares = [];
    const allowedLetters = ['r', 'n', 'b', 'k', 'q', 'p'];
    for (let i = 0; i < this.props.fen.length && this.props.fen[i] !== ' '; i++) {
      if (allowedLetters.includes(this.props.fen[i].toLowerCase())) {
        squares.push(this.props.fen.charAt(i));
      } else {
        if (this.props.fen.charAt(i) !== '/') {
          // its a number
          for(let j = 0; j < parseInt(this.props.fen[i], 10); j++) {
            squares.push('');
          }
        }
      }
    }
    let squareDivs = [];
    const colors = ['white', 'black'];
    let color = 0;
    let index = 0;
    for (let i=0; i <8; i++) {
      for (let j=0; j <8; j++) {
        squareDivs.push(<div id={index} key={i.toString() + j.toString()} style={{"display": "inline-block",  "width": "100px", "backgroundColor": colors[color], "height": "100px"}}>{squares[index]}</div>);
        color = (color + 1) % 2;
        index++;
      }
      color = (color + 1) % 2;
    }
    return (
      <div id="board" style={{'width': '800px'}}>
        { squareDivs }
      </div>
    );
  }
}