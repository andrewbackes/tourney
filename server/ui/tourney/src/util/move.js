export default class Move {
	static format(move) {
		if (move.source === 64) {
				return "-";
			}
			let sfile = String.fromCharCode(97 + (7 - (move.source % 8)));
			let srank = parseInt(move.source / 8, 10) + 1;
			let dfile = String.fromCharCode(97 + (7 - (move.destination % 8)));
			let drank = parseInt(move.destination / 8, 10) + 1;
			let alg = sfile + srank + dfile + drank;
			if (move.promote) {
				alg += move.promote;
			}
			return alg;
	}
}