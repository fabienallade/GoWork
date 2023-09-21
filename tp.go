package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tealeg/xlsx/v3"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// Task Creation de l'object qui doit
// contenir les différentes informations
type Task struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	Status      bool      `json:"status"`
}

func main() {
	// Creation de la variable pour choisir les options pour les differentes fonctions dans le terminal
	option := 10
	// Creation d'un tableau dynamique de type Task
	tasks := make([]Task, 0, 5)
	// Remplir la table des données venant du fichier JSON
	tasks = fromFile()

	for {
		switch option {
		case 0:
			fmt.Println()
			option = showMenu()
			break
		case 1:
			tasks = addTask(tasks)
			updateFile(tasks)
			option = 0
			break
		case 2:
			showTask(tasks)
			option = 0
			break
		case 3:
			tasks = markAsDone(tasks)
			updateFile(tasks)
			option = 0
			break
		case 4:
			tasks = deleteTask(tasks)
			updateFile(tasks)
			option = 0
			break
		case 5:
			exportXlsx(tasks)
			option = 0
			break
		case 6:
			fmt.Println("Au revoir")
			os.Exit(1)
		default:
			fmt.Println("Bienvenue dans mon application ALLADE FABIEN")
			option = showMenu()
			break
		}
	}
}

// Function d'affichage du menu
func showMenu() int {
	option := 0
	fmt.Println("1. Ajouter une tâche")
	fmt.Println("2. Afficher les tâches")
	fmt.Println("3. Marqué une tâche comme terminé")
	fmt.Println("4. Supprimer une tache")
	fmt.Println("5. Exporter le fichier en excel")
	fmt.Println("6. Quitter l'application")

	fmt.Println()
	fmt.Println("Veuillez choisir une option :")

	fmt.Scanf("%d \n", &option)

	return option
}

// Function pour afficher les taches depuis le fichier json
func showTask(tasks []Task) {
	fmt.Println("Vous avez choisi d'afficher les informtaion")
	showTables(tasks)
}

func addTask(tasks []Task) []Task {
	fmt.Println("Vous avez choisi l'ajout d'une nouvelle tache ")
	fmt.Println("Veuillez entrer le titre de la tâche")

	var title string
	fmt.Scanf("%s \n", &title)

	fmt.Println("Veuillez entrer la description de la task")
	var description string
	fmt.Scanf("%s \n", &description)

	// affecter la nouvelle taches
	task := Task{Title: title, Description: description, Status: false, CreatedAt: time.Now()}

	// Ajouter la dite taches a la variables taches
	newTasks := append(tasks, task)

	return newTasks
}

func showIfDone(status bool) string {
	if status {
		return "Terminer"
	} else {
		return "A faire"
	}
}

// Recuper les informations des fichiers JSON
func fromFile() []Task {
	// lecture du fichier json
	file, _ := ioutil.ReadFile("data/data.json")

	var data []Task

	// affecter les informations dans la variable data
	_ = json.Unmarshal([]byte(file), &data)

	return data
}

// Mettre a jour les fichiers
func updateFile(tasks []Task) {
	// ouverture du fichier d'écriture
	file, _ := json.MarshalIndent(tasks, "", " ")
	// affecter la nouvelle valeur de taches dans le fichier
	_ = ioutil.WriteFile("data/data.json", file, 0644)
}

// Marquer les taches comme etant effectuer
func markAsDone(tasks []Task) []Task {
	fmt.Println("Vous avez choisi l'option de marker un tache comme terminer")
	// afficher le tableau
	showTables(tasks)

	// Recuperer les informations du
	var choice int
	fmt.Println()
	fmt.Println("Veuillez choisir l'id de ce que vous voulez marquer comme terminer")

	fmt.Scanf("%d \n", &choice)

	// Mettre a jour le status du task choisi

	task := tasks[choice]

	task.Status = true

	tasks[choice] = task

	return tasks
}

// Supprimer une tache
func deleteTask(tasks []Task) []Task {
	fmt.Println("Vous avez choisi l'option de supprimer une tache ")
	showTables(tasks)

	var choice int
	fmt.Println()
	fmt.Println("Veuillez choisir l'id de ce que vous voulez supprimer")

	fmt.Scanf("%d \n", &choice)
	var newTasks []Task
	// Verifier si le choix qui a été fait existe dans les taches contenu dans la base de donnée
	if choice < len(tasks) {
		newTasks = append(tasks[:choice], tasks[choice+1:]...)
		fmt.Println("Suppression fait avec success")
	} else {
		fmt.Println("Vous ne pouvez pas supprimer un element qui n'existe pas")
	}

	return newTasks
}

// Exporter les informations des taches dans un fichier excel
func exportXlsx(tasks []Task) {
	fmt.Println("Vous avez choisi d'exorter le tableau")
	showTables(tasks)

	fmt.Println()
	fmt.Println("Veuillez mettre le nom du fichier que vous voulez generer")

	var name string
	fmt.Scanf("%s \n", &name)

	// Creer la configuration du nouveau fichier
	wb := xlsx.NewFile()
	// Ajouter la sheet Taches dans le nouveau fichiers crées
	sheet, _ := wb.AddSheet("Taches")

	// Ajouter un style pour les colones
	headerStyle := xlsx.NewStyle()
	headerStyle.Fill.FgColor = "FF00FF40"
	//Ajout de la couleur
	headerStyle.Fill.PatternType = "solid"
	// Appliquer la couleur
	headerStyle.ApplyFill = true

	// Definition des données des entetes du tableau
	headerList := []string{"ID", "STATUS", "TITRE", "DESCRIPTION", "CREE LE"}

	// Affecter les rows et le style et les valeurs
	headColumn := sheet.AddRow()
	for _, value := range headerList {
		row := headColumn.AddCell()
		row.SetStyle(headerStyle)
		row.Value = value
	}
	// Definir la taille des colomnes
	sheet.SetColWidth(1, 5, 15)

	// Ajouter les informations du json dans les rows du fichiers JSON
	for i, task := range tasks {
		row := sheet.AddRow()

		id := row.AddCell()
		id.Value = strconv.FormatInt(int64(i), 10)

		status := row.AddCell()
		status.Value = showIfDone(task.Status)

		title := row.AddCell()
		title.Value = task.Title

		description := row.AddCell()
		description.Value = task.Description

		createdAt := row.AddCell()
		createdAt.Value = task.CreatedAt.Format("2006-01-02 15:04:05")

	}

	// Enregistrer le nouveau fichier avec le noms qui a été reussi
	wb.Save("data/" + name + ".xlsx")
}

// Afficher la table avec le plugin
func showTables(tasks []Task) {
	table := tablewriter.NewWriter(os.Stdout)
	// Affecter le header
	table.SetHeader([]string{"ID", "STATUS", "TITRE", "DESCRIPTION", "DATE DE CREATION"})

	for i, v := range tasks {
		row := []string{strconv.FormatInt(int64(i), 10), showIfDone(v.Status),
			v.Title, v.Description, v.CreatedAt.Format("2006-01-02 15:04:05")}
		if v.Status {
			// Changer la couleur en fonction du status de la taches
			table.Rich(row, []tablewriter.Colors{tablewriter.Colors{tablewriter.Normal, tablewriter.BgHiGreenColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.BgHiGreenColor},
				tablewriter.Colors{tablewriter.Normal, tablewriter.BgHiGreenColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.BgHiGreenColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.BgHiGreenColor}})
		} else {
			table.Rich(row, []tablewriter.Colors{tablewriter.Colors{tablewriter.Normal, tablewriter.BgHiRedColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.BgHiRedColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.BgHiRedColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.BgHiRedColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.BgHiRedColor}})
		}

	}
	// Afficher les lignes ou plutot la bordures dans le tableau
	table.SetRowLine(true)

	// Afficher le tableau
	table.Render()
}
