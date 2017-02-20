package api

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Admins les administrateurs
type Admins struct {
	gorm.Model
	Email      string
	PassHash   string
	TokenReset string
	TokenLogin string
}

// Seasons Les saisons
type Seasons struct {
	gorm.Model
	Years string `gorm:"size:10"`
}

// Sports les sports
type Sports struct {
	gorm.Model
	Name string
}

// Categories les categories
type Categories struct {
	gorm.Model
	Name string
}

// Teams les équipes
type Teams struct {
	gorm.Model
	Name       string
	City       string
	Sport      Sports
	SportID    int
	Category   Categories
	CategoryID int
	Coaches    []Coaches `gorm:"many2many:coach_team;"`
	Players    []Players `gorm:"many2many:players_team;"`
}

// Players les joueurs
type Players struct {
	gorm.Model
	Number       int
	Fname        string
	Lname        string
	Email        string
	PassHash     string
	TokenRequest string
	TokenReset   string
	TokenLogin   string
	Teams        []Teams `gorm:"many2many:joueur_equipe;"`
}

// Locations les lieux
type Locations struct {
	gorm.Model
	Name    string
	City    string
	Address string
}

// Videos les videos
type Videos struct {
	gorm.Model
	Path      string
	Completed bool
}

// Games les parties
type Games struct {
	gorm.Model
	HomeTeam        Teams
	HomeTeamID      int
	OpposingTeam    Teams
	EquipeAdverseID int
	Season          Seasons
	SeasonID        int
	Location        Locations
	LocationID      int
	Video           Videos
	VideoID         int
	Date            string
	Action          []Actions
}

// Positions les positions des joueurs
type Positions struct {
	gorm.Model
	Name string
}

// PlayerPositionGameTeam table de relations entre
// les joueurs, les positions, les parties et les équipes
type PlayerPositionGameTeam struct {
	gorm.Model
	Player     Players
	PlayerID   int
	Game       Games
	GameID     int
	Position   Positions
	PositionID int
	Team       Teams
	TeamID     int
}

// Zones les zones de terrain
type Zones struct {
	gorm.Model
	Name string
}

//MovementsType represent Movement type entity
//1: Offensive
//2: Defensive
//3: Neutral
type MovementsType struct {
	gorm.Model
	Name string `gorm:"unique"`
}

// ActionsType les types d'actions
type ActionsType struct {
	gorm.Model
	Name           string
	Description    string
	MovementType   MovementsType
	MovementTypeID int
}

// Actions est une modélisation des informations sur une
// action exécutée par un joueur
type Actions struct {
	gorm.Model
	ActionType   ActionsType
	ActionTypeID int
	IsPositive   bool
	Zone         Zones
	ZoneID       int
	Game         Games
	GameID       int
	X1           float64
	Y1           float64
	X2           float64
	Y2           float64
	Time         time.Duration
	HomeScore    int
	GuestScore   int
	PLayer       Players
	PlayerID     int
}

// PlayersTeam table de relations joueurs et équipe
type PlayersTeam struct {
	Player   Players
	PlayerID int
	Team     Teams
	TeamID   int
}

// Coaches les entraineurs
type Coaches struct {
	gorm.Model
	Fname        string
	Lname        string
	Email        string
	PassHash     string
	TokenRequest string
	TokenReset   string
	TokenLogin   string
	Teams        []Teams `gorm:"many2many:entraineur_equipe;"`
}

// CoachTeam table de relations entre entraineurs et équipes
type CoachTeam struct {
	gorm.Model
	Coach   Coaches
	CoachID int
	Team    Teams
	TeamID  int
}

// Metrics les métriques
type Metrics struct {
	Name     string
	Equation string
	Team     Teams
	TeamID   int
}

// Expand effectue un fetch de tous les children de l'action
// (has-many, has-one, pas belongs-to)
func (a *Actions) Expand(db *gorm.DB) {
	db.Model(a).Related(&(a.ActionsType))
	db.Model(a).Related(&(a.Zones))
}

// Expand effectue un fetch de tous les children de la partie
// (has-many, has-one, pas belongs-to)
func (g *Games) Expand(db *gorm.DB) {
	db.Model(g).Related(&(g.Actions))
	db.Model(g).Related(&(g.HomeTeam), "HomeTeamID")
	db.Model(g).Related(&(g.OpposingTeam), "OpposingTeamID")
}
