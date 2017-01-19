package api

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TypeAction struct {
	gorm.Model
	Nom         string `gorm:"unique"`
	Description string
}
type Sport struct {
	gorm.Model
	Nom string
}
type Niveau struct {
	gorm.Model
	Nom string
}
type Entraineur struct {
	gorm.Model
	Prenom   string
	Nom      string
	Email    string
	PassHash string
	Token    string
	Equipes  []Equipe `gorm:"many2many:entraineur_equipe;"`
}
type Joueur struct {
	gorm.Model
	Prenom                string
	Nom                   string
	Numero                int
	Email                 string
	PassHash              string
	TokenInvitation       string
	TokenReinitialisation string
	TokenConnexion        string
	Equipes               []Equipe `gorm:"many2many:joueur_equipe;"`
}

type Equipe struct {
	gorm.Model
	Nom         string
	Ville       string
	Sport       Sport
	SportID     int
	Niveau      Niveau
	NiveauID    int
	Entraineurs []Entraineur `gorm:"many2many:entraineur_equipe;"`
	Joueurs     []Joueur     `gorm:"many2many:joueur_equipe;"`
}
type ActionZone struct {
	gorm.Model
	Nom string
}
type Saison struct {
	gorm.Model
	Annees string `gorm:"size:10"`
}
type Lieu struct {
	gorm.Model
	Nom     string
	Ville   string
	Adresse string
}
type Video struct {
	gorm.Model
	Path           string
	AnalyseTermine bool
}
type Zone struct {
	gorm.Model
	Nom string
}

type Partie struct {
	gorm.Model
	EquipeMaison    Equipe
	EquipeMaisonID  int
	EquipeAdverse   Equipe
	EquipeAdverseID int
	Saison          Saison
	SaisonID        int
	Lieu            Lieu
	LieuID          int
	Video           Video
	VideoID         int
	Date            string
	Actions         []Action
}

// Expand effectue un fetch de tous les children de la partie
// (has-many, has-one, pas belongs-to)
func (p *Partie) Expand(db *gorm.DB) {
	db.Model(p).Related(&(p.Actions))
	db.Model(p).Related(&(p.EquipeMaison), "EquipeMaisonID")
	db.Model(p).Related(&(p.EquipeAdverse), "EquipeAdverseID")
}

// Action est une modélisation des informations sur une
// action exécutée par un joueur
type Action struct {
	gorm.Model
	TypeAction      TypeAction
	TypeActionID    int
	ActionPositive  bool
	Zone            Zone
	ZoneID          int
	Partie          Partie
	PartieID        int
	X1              float64
	Y1              float64
	X2              float64
	Y2              float64
	Temps           time.Duration
	PointageMaison  int
	PointageAdverse int
	Joueur          Joueur
	JoueurID        int
}

// Expand effectue un fetch de tous les children de l'action
// (has-many, has-one, pas belongs-to)
func (a *Action) Expand(db *gorm.DB) {
	db.Model(a).Related(&(a.TypeAction))
	db.Model(a).Related(&(a.Zone))
}
