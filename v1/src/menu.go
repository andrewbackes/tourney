// TODO: this file is probably redundant with the command.go file. Should probably be deleted.

/*

 Project: Tourney

 Module: menu
 Description: Display menus.

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/15/2014

*/

package main

import (
	"fmt"
)

func menu() {
	/*

		This will prompt the user to answer questions about what type of
		tourney he/she wants to set up. This menus should appear when the
		program launches, unless

		The idea in designing this is to make it extremely easy to use.
		allow for multiple commands, even ones that dont answer the
		questions, so that the user wont get 'stuck' in a menu they dont
		want to answer

		allow for system wide commands that can be entered at any point

	*/

	fmt.Println("This is the main menu")
	var input string

	fmt.Scanf("%s", &input)
	fmt.Println("INPUT: ", input)
}
