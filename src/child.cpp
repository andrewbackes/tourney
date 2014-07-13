#include <iostream>
#include <string>
using namespace std;

int main() {
  char x;
  while(x != 'x') {
    cin >> x;
    if(x == 'a') {
      cout << 'A' << endl;
    } else {
      cout << 'B' << endl;
    }
  }
  return 0;
}