package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Daten struct {
	groesse      int
	baum         int
	monster      int //can be deleted if refractored to arraylength
	monster_posX []int
	monster_posY []int
	moves        []string
}

var daten Daten
var tbl [][]uint8
var gameEnd_bool bool

//===========================MAIN========================================

func main() {

	fmt.Printf("Starting server at port 8080\n")

	http.HandleFunc("/start", startHandler)

	http.HandleFunc("/game", inputHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)

	}

}

func initMoves() {
	daten.moves = append(daten.moves, "UL")
	daten.moves = append(daten.moves, "U")
	daten.moves = append(daten.moves, "UR")
	daten.moves = append(daten.moves, "R")
	daten.moves = append(daten.moves, "DR")
	daten.moves = append(daten.moves, "D")
	daten.moves = append(daten.moves, "DL")
	daten.moves = append(daten.moves, "L")
}

func startHandler(w http.ResponseWriter, r *http.Request) {

	// render Start Screen
	t, err := template.ParseFiles("start.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	t.Execute(w, nil)
	resetSession()

}

func inputHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	//GameEnd
	if gameEnd_bool {
		te, err := template.ParseFiles("gameover.html")
		if err != nil {
			fmt.Println(err.Error())

		}
		te.Execute(w, nil)
	} else {

		// Get Game Template
		t, err := template.ParseFiles("game.html")
		if err != nil {
			fmt.Println(err.Error())
		}

		input := r.FormValue("input")
		if input != "" {
			updateTbl(input)
			r.PostFormValue("input")
		}

		// Abfrage auf StartForm
		if r.FormValue("baum") != "" {
			//fmt.Fprintf(w, "POST request successful")
			groesse := r.FormValue("groesse")
			baum := r.FormValue("baum")
			monster := r.FormValue("monster")

			fmt.Printf("%s groesse - %s baum - %s monster", groesse, baum, monster)

			var err error

			daten.groesse, err = strconv.Atoi(groesse)
			if err != nil {
				fmt.Println(err.Error())
			}
			daten.baum, err = strconv.Atoi(baum)
			if err != nil {
				fmt.Println(err.Error())
			}
			daten.monster, err = strconv.Atoi(monster)
			if err != nil {
				fmt.Println(err.Error())
			}

			//return groesse, baum, monster

			//Abfrage Auf Gameinput
			//fmt.Println(getGridText())

		}
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, nil)
		fmt.Fprintf(w, "%s", getGridHtml(tbl))
		fmt.Println(getGridText(tbl))
	}
}

func newGrid() [][]uint8 {
	fmt.Println("Erstelle neues Spielfeld")
	initMoves()
	if daten.groesse <= 0 {
		fmt.Println("Error unzulässige Größe")
		return nil
	}
	tmp := Make2D[uint8](daten.groesse, daten.groesse)
	//daten.monster_posX = make([]int, daten.monster)
	//daten.monster_posY = make([]int, daten.monster)
	for z := daten.baum + daten.monster + 1; z > 0; z-- {
		x := rand.Intn(daten.groesse)
		y := rand.Intn(daten.groesse)

		if z > daten.baum+daten.monster && tmp[y][x] == 0 {
			// erstelle Spieler - id 3
			tmp[y][x] = 3
		} else if z > daten.monster && tmp[y][x] == 0 {
			// erstelle Bäume -id 1
			tmp[y][x] = 1
		} else if tmp[y][x] == 0 {
			// erstelle Monster -id 2
			tmp[y][x] = 2
			daten.monster_posX = append(daten.monster_posX, x)
			daten.monster_posY = append(daten.monster_posY, y)
		} else {
			z += 1
		}
		//fmt.Println(strconv.Itoa(x) + " - " + strconv.Itoa(y))
		//fmt.Println(getGridText(tmp))
	}
	return tmp
}
func getGridHtml(tabelle [][]uint8) string {
	if tabelle == nil {
		tabelle = newGrid()
	}

	var tbl_str string

	//tbl_str += "<div class = \"center\">"
	i := daten.groesse

	for x := 0; x < i; x++ {
		//Splate
		//tbl_str += "<div class=\"childCube\">"
		for y := 0; y < i; y++ {
			//Zeile
			tbl_str += "<div class=\"childCube\">"

			// print object in Grid---------------------

			switch tabelle[x][y] {
			case 0:
				{
					break
				}
			case 1:
				{
					// Baum Pic
					tbl_str += "B"
					break
				}
			case 2:
				{
					// Monster Pic
					tbl_str += "M"
					break
				}
			case 3:
				{
					// Player Pic
					tbl_str += "P"
					break
				}

			}
			tbl_str += "</div>"
			//-------------------------------------------------------

		}
		tbl_str += "<br><br>"
	}
	tbl = tabelle
	return tbl_str
}

func getGridText(tabelle [][]uint8) string {
	if tabelle == nil {
		tabelle = newGrid()
	}

	var tbl_str string

	//tbl_str += "<div id = \"1\">"
	i := daten.groesse
	// Ersetze durch DIV!
	for x := 0; x < i; x++ {
		//Splate
		tbl_str += "\n"
		//tbl_str += "[ ]"
		for y := 0; y < i; y++ {
			//Zeile

			// print object in Grid---------------------

			switch tabelle[x][y] {
			case 0:
				{
					tbl_str += "[ ]"
					break
				}
			case 1:
				{
					// Baum Pic
					tbl_str += "[B]"
					break
				}
			case 2:
				{
					// Monster Pic
					tbl_str += "[M]"
					break
				}
			case 3:
				{
					// Player Pic
					tbl_str += "[P]"
					break
				}

			}
			//-------------------------------------------------------

		}

	}
	tbl = tabelle
	return tbl_str
}

func updateTbl(input string) {
	//Update Player
	pos_y, pos_x := IndexOf(tbl, 3)
	if 0 < moveObj(pos_y, pos_x, input, -1) {
		gameEnd()
	}

	//Update Monster
	//TODO______________________________DEBUG_________________
	///*
	for m := 0; m < len(daten.monster_posX); m++ {
		r := rand.Intn(2) - 1

		var move string

		xComp := daten.monster_posX[m] - pos_x
		yComp := daten.monster_posY[m] - pos_y

		if xComp < 0 {
			// gehe rechts
			if yComp > 0 {
				// gehe hoch-rechts + benachbarte (OR)
				move = daten.moves[2+r]
			} else if yComp < 0 {
				// gehe runter-rechts ... (UR)
				move = daten.moves[4+r]

			} else if yComp == 0 {
				// nur rechts
				move = daten.moves[3+r]

			}

		} else if xComp > 0 {
			// gehe links

			if yComp > 0 {
				// gehe hoch-links + benachbarte
				// if -1 out of range fix
				if !(r == -1) {
					move = daten.moves[0+r]
				} else {
					move = daten.moves[7] //[len(daten.moves)]
				}
			} else if yComp < 0 {
				// gehe runter-links ...
				move = daten.moves[6+r]
			} else if yComp == 0 {
				// nur links
				// if -1 or 8 out of range fix
				if !(r == 1) {
					move = daten.moves[7+r]
				} else {
					move = daten.moves[0]
				}

			}

		} else if xComp == 0 {
			// Spieler Vertikal
			if yComp > 0 {
				// gehe nur hoch
				move = daten.moves[1+r]

			} else if yComp < 0 {
				// gehe nur runter
				move = daten.moves[5+r]
			}
		}

		//check collision with object
		coll := moveObj(daten.monster_posY[m],
			daten.monster_posX[m], move, m)
		//collision with tree or wall //-- nothing happens / monster dies
		if coll == 3 {
			//collision with player
			gameEnd()
		} else if coll != 0 {
			//delete monster
			tbl[daten.monster_posY[m]][daten.monster_posX[m]] = 0
			RemoveIndex(daten.monster_posX, m)
			RemoveIndex(daten.monster_posY, m)
			daten.monster -= 1

			if daten.monster == 0 {
				gameEnd()
			}
		}

	}
	//*/
}

// Go does not have optional parameters -.-
func moveObj(pos_x int, pos_y int, input string, m int) uint8 {

	obj := tbl[pos_x][pos_y]
	new_posX, new_posY := -1, -1

	switch input {

	case "UL":
		{
			if checkCol(pos_x-1, pos_y-1) == 0 {
				new_posX = pos_x - 1
				new_posY = pos_y - 1

			} else {
				return checkCol(pos_x-1, pos_y-1)
			}
			break
		}
	case "U":
		{
			if checkCol(pos_x-1, pos_y) == 0 {
				new_posX = pos_x - 1
				new_posY = pos_y
			} else {
				return checkCol(pos_x-1, pos_y)
			}
			break
		}
	case "UR":
		{
			if checkCol(pos_x-1, pos_y+1) == 0 {
				new_posX = pos_x - 1
				new_posY = pos_y + 1
			} else {
				return checkCol(pos_x-1, pos_y+1)
			}
			break
		}
	case "L":
		{
			if checkCol(pos_x, pos_y-1) == 0 {
				new_posX = pos_x
				new_posY = pos_y - 1
			} else {
				return checkCol(pos_x, pos_y-1)
			}
			break
		}
	case "R":
		{
			if checkCol(pos_x, pos_y+1) == 0 {
				new_posX = pos_x
				new_posY = pos_y + 1
			} else {
				return checkCol(pos_x, pos_y+1)
			}
			break
		}
	case "DL":
		{
			if checkCol(pos_x+1, pos_y-1) == 0 {
				new_posX = pos_x + 1
				new_posY = pos_y - 1
			} else {
				return checkCol(pos_x+1, pos_y-1)
			}
			break
		}
	case "D":
		{

			if checkCol(pos_x+1, pos_y) == 0 {
				new_posX = pos_x + 1
				new_posY = pos_y
			} else {
				return checkCol(pos_x+1, pos_y)
			}
			break
		}
	case "DR":
		{
			if checkCol(pos_x+1, pos_y+1) == 0 {
				new_posX = pos_x + 1
				new_posY = pos_y + 1
			} else {
				return checkCol(pos_x+1, pos_y+1)
			}
			break
		}
	}
	//set old position and new
	if !(new_posX == -1 && new_posY == -1) {
		tbl[pos_x][pos_y] = 0
		tbl[new_posX][new_posY] = obj

		//update monster Array
		if m > -1 {
			daten.monster_posX[m] = new_posY // parameterswap in func start
			daten.monster_posY[m] = new_posX
		}
	}
	return 0
}

func gameEnd() {
	// Show Endscreen Leaderboard
	fmt.Println("GAME END!")
	gameEnd_bool = true

}
func Make2D[T any](n, m int) [][]T {
	matrix := make([][]T, n)
	rows := make([]T, n*m)
	for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
		endRow := startRow + m
		matrix[i] = rows[startRow:endRow:endRow]
	}
	return matrix
}

func IndexOf[P uint8](arr [][]P, str P) (int, int) {
	for x := range arr {
		for y, t := range arr[x] {
			if t == str {
				return x, y
			}
		}
	}
	return -1, -1
}
func checkCol(x int, y int) uint8 {
	if x < 0 || y < 0 || x > daten.groesse-1 || y > daten.groesse-1 {
		return 4
	} else {
		tmp := tbl[x][y]
		return tmp
	}
}
func resetSession() {
	tbl = nil
	daten = Daten{}
	gameEnd_bool = false
}

func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}
