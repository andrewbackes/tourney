#include <iostream>

int main() {
  char x;
  std::cout << "child2" << std::endl;
  for(std::cin >> x ; x != 'x' ; std::cin >> x) {
    if(x == 'a')
      std::cout << "Asdf" << std::endl;
    else
      std::cout << "Bsdf" << std::endl;
  }
  return 0;
}