import React, { Component } from 'react';
import WhitePawn from 'images/pieces/piece11.png';
import WhiteKnight from 'images/pieces/piece12.png';
import WhiteBishop from'images/pieces/piece13.png';
import WhiteRook from 'images/pieces/piece14.png';
import WhiteQueen from 'images/pieces/piece15.png';
import WhiteKing from 'images/pieces/piece16.png';
import BlackPawn from 'images/pieces/piece01.png';
import BlackKnight from 'images/pieces/piece02.png';
import BlackBishop from 'images/pieces/piece03.png';
import BlackRook from 'images/pieces/piece04.png';
import BlackQueen from 'images/pieces/piece05.png';
import BlackKing from 'images/pieces/piece06.png';

export default class Board extends Component {
  render() {
    let squares = [];
    const allowedLetters = ['r', 'n', 'b', 'k', 'q', 'p'];
    const images = {
      'R': <img alt="R" src={WhiteRook} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>,
      'N': <img alt="R" src={WhiteKnight} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>,
      'B': <img alt="R" src={WhiteBishop} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>,
      'Q': <img alt="R" src={WhiteQueen} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>,
      'K': <img alt="R" src={WhiteKing} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>,
      'P': <img alt="R" src={WhitePawn} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>,
      'r': <img alt="R" src={BlackRook} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>,
      'n': <img alt="R" src={BlackKnight} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>,
      'b': <img alt="R" src={BlackBishop} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>,
      'q': <img alt="R" src={BlackQueen} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>,
      'k': <img alt="R" src={BlackKing} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>,
      'p': <img alt="R" src={BlackPawn} style={{ "width": "30px", "height": "30px", "marginLeft": "10px", "marginTop": "10px" }}/>
    };
    for (let i = 0; i < this.props.position.fen.length && this.props.position.fen[i] !== ' '; i++) {
      if (allowedLetters.includes(this.props.position.fen[i].toLowerCase())) {
        squares.push(this.props.position.fen.charAt(i));
      } else {
        if (this.props.position.fen.charAt(i) !== '/') {
          // its a number
          for(let j = 0; j < parseInt(this.props.position.fen[i], 10); j++) {
            squares.push('');
          }
        }
      }
    }
    let squareDivs = [];
    const colors = ['white', 'gray'];
    let color = 0;
    let index = 0;
    for (let i=0; i <8; i++) {
      for (let j=0; j <8; j++) {
        let border = '1px black solid';
        if (this.props.position && this.props.position.lastMove) {
          if ((63 - index) === this.props.position.lastMove.source || (63 - index) === this.props.position.lastMove.destination) {
            border = '2px yellow solid';
          }
        }
        squareDivs.push(
          <div id={index} key={i.toString() + j.toString()} style={{ 
            "border": border,
            "display": "inline-block",
            "width": "50px",
            "height": "50px",
            "position": "absolute",
            "top": i * 50 + "px",
            "left": j * 50 + "px",
            "backgroundColor": colors[color]
          }}>
            {images[squares[index]]}
          </div>);
        color = (color + 1) % 2;
        index++;
      }
      color = (color + 1) % 2;
    }
    return (
      <div id="board" style={{ 'position':'relative', 'width': '400px', 'height': '400px'}}>
        { squareDivs }
      </div>
    );
  }
}