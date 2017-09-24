
export default class Duration {
  static format(duration) {
    //let hrUnit = 3600000000000; // nanoseconds
    //let minUnit = 60000000000; // nanoseconds
    //let secUnit = 1000000000; // nanoseconds
    //let msUnit = 1000000; // nanoseconds
    
    let pretty = "";
    
    let min = parseInt(duration / 60000000000, 10);
    if (min !== 0) { pretty += min + "m "; }
    let remainder = duration % 60000000000;

    let ms = parseInt(remainder / 1000000, 10);
    if (ms !== 0) { pretty += ms/1000 + "s"; }

    return pretty;
  }
}