import Duration from 'util/duration';

export default class TimeControl {
  static format(timeControl) {
    let val = "";
    val += (timeControl.moves ? timeControl.moves.toString() : "-") + "/";
    val += timeControl.time ? Duration.format(timeControl.time.toString()) : "-";
    val += timeControl.increment ? "+" + Duration.format(timeControl.increment.toString()) : "";
    return val;
  }
}