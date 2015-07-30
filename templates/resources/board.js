/* ********************************************************************************

	BOARD DRAWING AND FEN COMPUTING:

******************************************************************************** */			

	var pawn = new Array(2);
	pawn[0] = "<img src='resources/pieces/piece11.png' class='piece' style='width:24px;height:32px; position:relative;left:8px;top:4px;' />";
	pawn[1] = "<img src='resources/pieces/piece01.png' class='piece' style='width:24px;height:32px; position:relative;left:8px;top:4px;' />";
	var knight = new Array(2);
	knight[0] = "<img src='resources/pieces/piece12.png' class='piece' style='width:32px;height:32px; position:relative;left:3px;top:4px;'/>";
	knight[1] = "<img src='resources/pieces/piece02.png' class='piece' style='width:32px;height:32px; position:relative;left:3px;top:4px;'/>";
	var bishop = new Array(2);
	bishop[0] = "<img src='resources/pieces/piece13.png' class='piece' style='width:30px;height:32px; position:relative;left:3px;top:4px;' />";
	bishop[1] = "<img src='resources/pieces/piece03.png' class='piece' style='width:30px;height:32px; position:relative;left:3px;top:4px;' />";
	var rook = new Array(2);
	rook[0] = "<img src='resources/pieces/piece14.png' class='piece' style='width:27px;height:32px; position:relative;left:5px;top:4px;'/>";
	rook[1] = "<img src='resources/pieces/piece04.png' class='piece' style='width:27px;height:32px; position:relative;left:5px;top:4px;'/>";
	var queen = new Array(2);
	queen[0] = "<img src='resources/pieces/piece15.png' class='piece' style='width:32px;height:30px; position:relative;left:3px;top:4px;'/>";
	queen[1] = "<img src='resources/pieces/piece05.png' class='piece' style='width:32px;height:30px; position:relative;left:3px;top:4px;'/>";
	var king = new Array(2);
	king[0] = "<img src='resources/pieces/piece16.png' class='piece' style='width:32px;height:32px; position:relative;left:2px;top:4px;'/>";
	king[1] = "<img src='resources/pieces/piece06.png' class='piece' style='width:32px;height:32px; position:relative;left:2px;top:4px;'/>";

	function idFromIndex(i) {
		var rank = Math.floor(i/8)+1;
		var file = (Math.floor(i%8)+1);
		var idName = String.fromCharCode(96+file) + rank;
		return idName;
	}

	function clearBoard() {
		for(i=0;i<64;i++) {
			document.getElementById(idFromIndex(i)).innerHTML = "";
		}
	}

	function setNewBoard() {
		setBoard(startingFEN);
	}

	var pieces = {
		'P': pawn[0], 'p': pawn[1],
		'K': king[0], 'k': king[1],
		'R': rook[0], 'r': rook[1],
		'N': knight[0], 'n': knight[1],
		'B': bishop[0], 'b': bishop[1],
		'Q': queen[0], 'q': queen[1]
	};
	
	function setBoard(FEN) {
		var row = 7;
		var col = 0;
		for (var pos=0; FEN.charAt(pos) != " "; pos++) {
			var char = FEN.charAt(pos);
			if( char >= '0' && char <= '9') {
				for(var i=0; i < parseInt(char); i++) {
					var square = (8*row) + col+i;
					document.getElementById(idFromIndex(square)).innerHTML = "";
				}
				col+=parseInt(char);

			} else if (char == "/") {
				row--;
				col = 0;
			} else {
				var square = (8*row) + col;
				document.getElementById(idFromIndex(square)).innerHTML = pieces[char];
				col++;
			}
		}
	}