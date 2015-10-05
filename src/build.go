/*******************************************************************************

 Project: Tourney

 Module: build

 Description: build engines

 Author(s): Andrew Backes
 Created: 10/03/2015

*******************************************************************************/

package main

type BuildSpec struct {
	Repo string
	Branch string
	BuildFile string
}