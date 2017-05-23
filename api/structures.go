package api

import "github.com/jinzhu/gorm"

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
	Season     Seasons
	SeasonID   int
	Sexe       string
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
	Name         string
	City         string
	Address      string
	IsInside     bool
	FieldType    FieldTypes
	FieldTypesID int
}

// FieldTypes les types de terrains
type FieldTypes struct {
	gorm.Model
	Type        string
	Description string
}

// Videos les videos
type Videos struct {
	gorm.Model
	Path      string
	Part      int
	Completed int
	Game      Games
	GameID    int
}

// Games les parties
type Games struct {
	gorm.Model
	Team           Teams
	TeamID         int
	Status         string // Local/visiteur
	OpposingTeam   string
	Season         Seasons
	SeasonID       int
	Location       Locations
	LocationID     int
	FieldCondition string
	Date           string
	Temperature    string // Rain, wind, sun, etc
	Degree         string // The temperature in degree celcius
	Action         []Actions
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
	Description string
	TypeAction  string
	Name        string
}

// Actions est une modélisation des informations sur une
// action exécutée par un joueur
type Actions struct {
	gorm.Model
	ActionType      ActionsType
	ActionTypeID    int
	ReceptionType   ReceptionType
	ReceptionTypeID int
	Zone            Zones
	ZoneID          int
	Game            Games
	GameID          int
	X1              float64
	Y1              float64
	X2              float64
	Y2              float64
	X3              float64
	Y3              float64
	Time            float64
	HomeScore       int
	GuestScore      int
	Player          Players
	PlayerID        int
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
	CoachID      int
	Fname        string
	Lname        string
	Email        string
	PassHash     string
	TokenRequest string
	TokenReset   string
	TokenLogin   string
	Teams        []Teams `gorm:"many2many:coach_team;"`
}

// CoachTeam table de relations entre entraineurs et équipes
type CoachTeam struct {
	gorm.Model
	Coach   Coaches
	CoachID int
	Team    Teams
	TeamID  int
}

// ReceptionType Reception du ballon
type ReceptionType struct {
	gorm.Model
	Name string
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
/*func (a *Actions) Expand(db *gorm.DB) {
	db.Model(a).Related(&(a.ActionType))
	db.Model(a).Related(&(a.Zone))
}

// Expand effectue un fetch de tous les children de la partie
// (has-many, has-one, pas belongs-to)
func (g *Games) Expand(db *gorm.DB) {
	db.Model(g).Related(&(g.Action))
	db.Model(g).Related(&(g.Team), "TeamID")
	db.Model(g).Related(&(g.OpposingTeam), "OpposingTeamID")
}*/
