package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// Game constants
const (
	defaultWidth          = 800
	defaultHeight         = 600
	gravity               = 0.5
	jumpForce             = -12
	doubleJumpForce       = -10
	maxDoubleJumps        = 1
	playerSize            = 30
	baseSpeed             = 3.0
	speedIncrease         = 0.2
	speedIncreaseInterval = 20
	startPlatformDuration = 180
	minPlatformDistance   = 150.0
	maxPlatformDistance   = 250.0
	platformFadeDistance  = 200.0
)

// Game state constants
const (
	stateMenu = iota
	stateSettings
	statePlaying
	stateGameOver
	statePause
	stateNameInput
	stateLeaderboard
)

// Resolution options
var availableResolutions = []struct {
	w, h int
	name string
}{
	{800, 600, "800x600"},
	{1280, 720, "1280x720"},
	{1920, 1080, "1920x1080"},
	{2560, 1440, "2560x1440"},
}

// Platform levels with colors
var platformLevels = []struct {
	y     float64
	color color.Color
	name  string
}{
	{float64(defaultHeight) - 50, color.RGBA{0, 255, 0, 255}, "GREEN"},
	{float64(defaultHeight) - 120, color.RGBA{0, 200, 255, 255}, "BLUE"},
	{float64(defaultHeight) - 190, color.RGBA{255, 255, 0, 255}, "YELLOW"},
	{float64(defaultHeight) - 260, color.RGBA{255, 165, 0, 255}, "ORANGE"},
	{float64(defaultHeight) - 330, color.RGBA{255, 0, 0, 255}, "RED"},
}

type Player struct {
	x, y          float64
	velX, velY    float64
	isJumping     bool
	color         color.Color
	doubleJumps   int
	canDoubleJump bool
}

type Platform struct {
	x, y, width float64
	color       color.Color
	level       int
	alpha       float64
}

type HighScore struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type Game struct {
	player        Player
	platforms     []Platform
	cameraX       float64
	score         int
	lastPlatformX float64
	startTimer    int
	gameState     int
	currentSpeed  float64
	jumpKey       ebiten.Key
	screenWidth   int
	screenHeight  int
	settingIndex  int
	resolutionIdx int
	nameInput     string
	highScores    []HighScore
	fontScale     float64
}

func (g *Game) Update() error {
	switch g.gameState {
	case stateMenu:
		return g.updateMenu()
	case stateSettings:
		return g.updateSettings()
	case statePlaying:
		return g.updateGame()
	case stateGameOver:
		return g.updateGameOver()
	case statePause:
		return g.updatePause()
	case stateNameInput:
		return g.updateNameInput()
	case stateLeaderboard:
		return g.updateLeaderboard()
	}
	return nil
}

func (g *Game) updateMenu() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.settingIndex = (g.settingIndex + 1) % 4
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.settingIndex = (g.settingIndex + 3) % 4
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch g.settingIndex {
		case 0:
			g.resetGame()
			g.gameState = statePlaying
		case 1:
			g.gameState = stateLeaderboard
		case 2:
			g.gameState = stateSettings
			g.settingIndex = 0
		case 3:
			return ebiten.Termination
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return nil
}

func (g *Game) updateLeaderboard() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.gameState = stateMenu
	}
	return nil
}

func (g *Game) updateSettings() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.settingIndex = (g.settingIndex + 1) % 3
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.settingIndex = (g.settingIndex + 2) % 3
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch g.settingIndex {
		case 0:
			g.resolutionIdx = (g.resolutionIdx + 1) % len(availableResolutions)
			g.screenWidth = availableResolutions[g.resolutionIdx].w
			g.screenHeight = availableResolutions[g.resolutionIdx].h
			g.updateFontScale()
			ebiten.SetWindowSize(g.screenWidth, g.screenHeight)
		case 1:
			keys := inpututil.AppendPressedKeys(nil)
			if len(keys) > 0 {
				g.jumpKey = keys[0]
			}
		case 2:
			g.gameState = stateMenu
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.gameState = stateMenu
	}
	return nil
}

func (g *Game) updateFontScale() {
	baseResolution := 800.0
	g.fontScale = float64(g.screenWidth) / baseResolution
}

func (g *Game) updateGame() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.gameState = statePause
		return nil
	}

	if g.startTimer > 0 {
		g.startTimer--
		if inpututil.IsKeyJustPressed(g.jumpKey) {
			g.startTimer = 0
		}
		return nil
	}

	if g.score > 0 && g.score%speedIncreaseInterval == 0 {
		g.currentSpeed = baseSpeed * (1 + speedIncrease*float64(g.score/speedIncreaseInterval))
	}

	if inpututil.IsKeyJustPressed(g.jumpKey) {
		if !g.player.isJumping {
			g.player.velY = jumpForce
			g.player.isJumping = true
			g.player.canDoubleJump = true
			g.player.doubleJumps = maxDoubleJumps
		} else if g.player.canDoubleJump && g.player.doubleJumps > 0 {
			g.player.velY = doubleJumpForce
			g.player.doubleJumps--
			if g.player.doubleJumps == 0 {
				g.player.canDoubleJump = false
			}
		}
	}

	g.player.velY += gravity
	g.player.x += g.currentSpeed
	g.player.y += g.player.velY

	g.handlePlatformCollision()

	if g.player.y > float64(g.screenHeight) {
		g.gameState = stateNameInput
		return nil
	}

	g.generatePlatforms()
	g.updatePlatformFade()
	g.cleanPassedPlatforms()

	g.cameraX = g.player.x - float64(g.screenWidth)*0.4
	return nil
}

func (g *Game) updatePlatformFade() {
	for i := range g.platforms {
		dist := g.player.x - (g.platforms[i].x + g.platforms[i].width)
		if dist > 0 {
			g.platforms[i].alpha = 1.0 - math.Min(dist/platformFadeDistance, 1.0)
		} else {
			g.platforms[i].alpha = 1.0
		}
	}
}

func (g *Game) handlePlatformCollision() {
	g.player.isJumping = true
	for _, p := range g.platforms {
		if g.player.y+playerSize >= p.y &&
			g.player.y+playerSize <= p.y+20 &&
			g.player.x+playerSize > p.x &&
			g.player.x < p.x+p.width &&
			g.player.velY >= 0 {

			g.player.y = p.y - playerSize
			g.player.velY = 0
			g.player.isJumping = false
			g.player.doubleJumps = maxDoubleJumps
			g.player.canDoubleJump = true
			break
		}
	}
}

func (g *Game) generatePlatforms() {
	if g.lastPlatformX > g.player.x+float64(g.screenWidth)-200 {
		return
	}

	tUp := -jumpForce / gravity
	maxJumpDistance := g.currentSpeed * (2 * tUp)
	safeDistance := math.Min(maxJumpDistance*0.8, maxPlatformDistance)

	currentLevel := 0
	minDist := math.MaxFloat64
	for i, level := range platformLevels {
		dist := math.Abs(g.player.y - (level.y - playerSize))
		if dist < minDist {
			minDist = dist
			currentLevel = i
		}
	}

	nextLevel := currentLevel
	r := rand.Intn(3)
	if r == 0 && currentLevel > 0 {
		nextLevel = currentLevel - 1
	} else if r == 2 && currentLevel < len(platformLevels)-1 {
		nextLevel = currentLevel + 1
	}

	newX := g.lastPlatformX + minPlatformDistance + rand.Float64()*(safeDistance-minPlatformDistance)
	newWidth := 80.0 + rand.Float64()*70.0

	g.platforms = append(g.platforms, Platform{
		x:     newX,
		y:     platformLevels[nextLevel].y,
		width: newWidth,
		color: platformLevels[nextLevel].color,
		level: nextLevel,
		alpha: 1.0,
	})

	g.lastPlatformX = newX
}

func (g *Game) cleanPassedPlatforms() {
	for i := 0; i < len(g.platforms); {
		if g.platforms[i].x+g.platforms[i].width < g.cameraX-100 {
			g.platforms = append(g.platforms[:i], g.platforms[i+1:]...)
			g.score++
		} else {
			i++
		}
	}
}

func (g *Game) updateGameOver() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.resetGame()
		g.gameState = statePlaying
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.gameState = stateMenu
	}
	return nil
}

func (g *Game) updatePause() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.gameState = statePlaying
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.gameState = stateMenu
	}
	return nil
}

func (g *Game) updateNameInput() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(g.nameInput) > 0 {
			g.nameInput = g.nameInput[:len(g.nameInput)-1]
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if len(g.nameInput) > 0 {
			g.addHighScore(g.nameInput, g.score)
			g.gameState = stateGameOver
		}
		return nil
	}

	g.nameInput += string(ebiten.InputChars())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 40, 255})

	switch g.gameState {
	case stateMenu:
		g.drawMenu(screen)
	case stateSettings:
		g.drawSettings(screen)
	case statePlaying:
		g.drawGame(screen)
	case stateGameOver:
		g.drawGameOver(screen)
	case statePause:
		g.drawGame(screen)
		g.drawPause(screen)
	case stateNameInput:
		g.drawNameInput(screen)
	case stateLeaderboard:
		g.drawLeaderboard(screen)
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	title := "GEOMETRY DASH"
	options := []string{"START GAME", "LEADERBOARD", "SETTINGS", "EXIT"}

	white := color.White
	yellow := color.RGBA{255, 255, 0, 255}

	titleX := g.screenWidth/2 - 70
	titleY := g.screenHeight / 3
	optionStartX := g.screenWidth / 3
	optionStartY := g.screenHeight / 2
	optionSpacing := 30

	text.Draw(screen, title, basicfont.Face7x13, titleX, titleY, white)

	for i, opt := range options {
		var clr color.Color = white
		if i == g.settingIndex {
			clr = yellow
			text.Draw(screen, ">", basicfont.Face7x13,
				optionStartX,
				optionStartY+i*optionSpacing,
				clr)
		}
		text.Draw(screen, opt, basicfont.Face7x13,
			optionStartX+20,
			optionStartY+i*optionSpacing,
			clr)
	}

	instructionsY := g.screenHeight - 50
	text.Draw(screen, "USE ARROW KEYS TO NAVIGATE", basicfont.Face7x13,
		g.screenWidth/4, instructionsY, white)
	text.Draw(screen, "PRESS ENTER TO SELECT", basicfont.Face7x13,
		g.screenWidth/4, instructionsY+20, white)
}

func (g *Game) drawLeaderboard(screen *ebiten.Image) {
	title := "LEADERBOARD"
	white := color.White

	titleX := g.screenWidth/2 - 60
	titleY := g.screenHeight / 4
	listStartX := g.screenWidth/2 - 100
	listStartY := g.screenHeight / 3
	itemSpacing := 30
	returnY := g.screenHeight - 50

	text.Draw(screen, title, basicfont.Face7x13, titleX, titleY, white)

	if len(g.highScores) == 0 {
		text.Draw(screen, "NO SCORES YET", basicfont.Face7x13,
			g.screenWidth/2-50,
			g.screenHeight/2, white)
	} else {
		for i, hs := range g.highScores {
			if i >= 10 {
				break
			}
			text.Draw(screen, fmt.Sprintf("%d. %s: %d", i+1, hs.Name, hs.Score),
				basicfont.Face7x13,
				listStartX,
				listStartY+i*itemSpacing,
				white)
		}
	}

	text.Draw(screen, "PRESS ESC TO RETURN", basicfont.Face7x13,
		g.screenWidth/2-80,
		returnY, white)
}

func (g *Game) drawSettings(screen *ebiten.Image) {
	title := "SETTINGS"
	settings := []string{
		fmt.Sprintf("RESOLUTION: %s", availableResolutions[g.resolutionIdx].name),
		fmt.Sprintf("JUMP KEY: %s", g.jumpKey.String()),
		"BACK TO MENU",
	}

	white := color.White
	yellow := color.RGBA{255, 255, 0, 255}

	text.Draw(screen, title, basicfont.Face7x13,
		g.screenWidth/2-40,
		g.screenHeight/4, white)

	for i, setting := range settings {
		var currentColor color.Color = white
		if i == g.settingIndex {
			currentColor = yellow
			text.Draw(screen, ">", basicfont.Face7x13,
				g.screenWidth/3,
				g.screenHeight/3+i*30,
				currentColor)
		}
		text.Draw(screen, setting, basicfont.Face7x13,
			g.screenWidth/3+30,
			g.screenHeight/3+i*30,
			currentColor)
	}

	text.Draw(screen, "USE ARROW KEYS TO CHANGE SETTINGS",
		basicfont.Face7x13, g.screenWidth/4, g.screenHeight-50, white)
	text.Draw(screen, "PRESS ENTER TO CONFIRM",
		basicfont.Face7x13, g.screenWidth/4, g.screenHeight-30, white)
}

func (g *Game) drawGame(screen *ebiten.Image) {
	for _, p := range g.platforms {
		col := p.color.(color.RGBA)
		col.A = uint8(p.alpha * 255)
		ebitenutil.DrawRect(screen,
			p.x-g.cameraX,
			p.y,
			p.width, 20,
			col)
	}

	ebitenutil.DrawRect(screen,
		g.player.x-g.cameraX,
		g.player.y,
		playerSize, playerSize,
		g.player.color)

	if g.startTimer <= 0 {
		text.Draw(screen,
			fmt.Sprintf("SCORE: %d", g.score),
			basicfont.Face7x13,
			20,
			20,
			color.White)
		text.Draw(screen,
			fmt.Sprintf("SPEED: %.1fx", g.currentSpeed/baseSpeed),
			basicfont.Face7x13,
			20,
			40,
			color.White)
		text.Draw(screen,
			fmt.Sprintf("JUMPS: %d", g.player.doubleJumps),
			basicfont.Face7x13,
			20,
			60,
			color.White)
	} else {
		text.Draw(screen,
			"PRESS JUMP TO START",
			basicfont.Face7x13,
			g.screenWidth/2-70,
			g.screenHeight/2,
			color.White)
		text.Draw(screen,
			fmt.Sprintf("(USE %s KEY)", g.jumpKey.String()),
			basicfont.Face7x13,
			g.screenWidth/2-50,
			g.screenHeight/2+20,
			color.White)
	}
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	white := color.White

	centerX := g.screenWidth / 2
	gameOverY := g.screenHeight / 4
	scoreY := gameOverY + 30
	highScoresTitleY := g.screenHeight / 2
	restartY := g.screenHeight - 80
	menuY := g.screenHeight - 50

	text.Draw(screen, "GAME OVER", basicfont.Face7x13,
		centerX-40, gameOverY, white)

	text.Draw(screen, fmt.Sprintf("YOUR SCORE: %d", g.score),
		basicfont.Face7x13, centerX-50, scoreY, white)

	text.Draw(screen, "HIGH SCORES:", basicfont.Face7x13,
		centerX-50, highScoresTitleY, white)

	for i, hs := range g.highScores {
		if i >= 5 {
			break
		}
		text.Draw(screen, fmt.Sprintf("%d. %s: %d", i+1, hs.Name, hs.Score),
			basicfont.Face7x13, centerX-50, highScoresTitleY+20+i*20, white)
	}

	text.Draw(screen, "PRESS R TO RESTART", basicfont.Face7x13,
		centerX-70, restartY, white)
	text.Draw(screen, "PRESS ESC FOR MENU", basicfont.Face7x13,
		centerX-70, menuY, white)
}

func (g *Game) drawPause(screen *ebiten.Image) {
	text.Draw(screen,
		"PAUSE",
		basicfont.Face7x13,
		g.screenWidth/2-30,
		g.screenHeight/2-40,
		color.White)
	text.Draw(screen,
		"Press ESC to continue",
		basicfont.Face7x13,
		g.screenWidth/2-70,
		g.screenHeight/2,
		color.White)
	text.Draw(screen,
		"Press M for menu",
		basicfont.Face7x13,
		g.screenWidth/2-60,
		g.screenHeight/2+30,
		color.White)
}

func (g *Game) drawNameInput(screen *ebiten.Image) {
	text.Draw(screen,
		"ENTER YOUR NAME:",
		basicfont.Face7x13,
		g.screenWidth/2-70,
		g.screenHeight/2-20,
		color.White)
	text.Draw(screen,
		g.nameInput,
		basicfont.Face7x13,
		g.screenWidth/2-50,
		g.screenHeight/2+10,
		color.White)
	text.Draw(screen,
		"PRESS ENTER TO CONFIRM",
		basicfont.Face7x13,
		g.screenWidth/2-90,
		g.screenHeight/2+40,
		color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

func (g *Game) resetGame() {
	g.player = Player{
		x:     float64(g.screenWidth) * 0.4,
		y:     platformLevels[0].y - playerSize,
		velX:  0,
		velY:  0,
		color: color.RGBA{255, 100, 100, 255},
	}

	g.platforms = []Platform{
		{
			x:     g.player.x - 100,
			y:     platformLevels[0].y,
			width: 200,
			color: platformLevels[0].color,
			level: 0,
			alpha: 1.0,
		},
		{
			x:     g.player.x + 150,
			y:     platformLevels[1].y,
			width: 100,
			color: platformLevels[1].color,
			level: 1,
			alpha: 1.0,
		},
	}

	g.cameraX = 0
	g.score = 0
	g.lastPlatformX = g.player.x + 150
	g.startTimer = startPlatformDuration
	g.currentSpeed = baseSpeed
	g.player.doubleJumps = maxDoubleJumps
	g.player.canDoubleJump = true
}

func (g *Game) loadHighScores() error {
	file, err := os.ReadFile("highscores.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(file, &g.highScores)
}

func (g *Game) saveHighScores() error {
	data, err := json.Marshal(g.highScores)
	if err != nil {
		return err
	}
	return os.WriteFile("highscores.json", data, 0644)
}

func (g *Game) addHighScore(name string, score int) {
	g.highScores = append(g.highScores, HighScore{Name: name, Score: score})

	for i := 0; i < len(g.highScores); i++ {
		for j := i + 1; j < len(g.highScores); j++ {
			if g.highScores[i].Score < g.highScores[j].Score {
				g.highScores[i], g.highScores[j] = g.highScores[j], g.highScores[i]
			}
		}
	}

	if len(g.highScores) > 10 {
		g.highScores = g.highScores[:10]
	}

	g.saveHighScores()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		gameState:     stateMenu,
		jumpKey:       ebiten.KeySpace,
		screenWidth:   availableResolutions[0].w,
		screenHeight:  availableResolutions[0].h,
		resolutionIdx: 0,
		currentSpeed:  baseSpeed,
		nameInput:     "Player",
		fontScale:     1.0,
	}

	game.updateFontScale()

	if err := game.loadHighScores(); err != nil {
		log.Printf("Error loading high scores: %v", err)
	}

	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("Geometry Dash Clone")
	ebiten.SetWindowResizable(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
