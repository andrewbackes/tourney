#include <iostream>

int main() {
  std::cout << "child1" << std::endl;
  char x;
  for(std::cin >> x ; x != 'x' ; std::cin >> x) {
    if(x == 'a') std::cout << "A" << std::endl;
    else std::cout << "B" << std::endl;
  }
  return 0;
}