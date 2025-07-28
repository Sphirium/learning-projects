package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	defaultWidth          = 800
	defaultHeight         = 600
	gravity               = 0.5
	jumpForce             = -12
	playerSize            = 30
	baseSpeed             = 3.0
	speedIncrease         = 0.2
	speedIncreaseInterval = 30
	startPlatformDuration = 180
	slopeLength           = 300.0
	slopeHeight           = 50.0
)

var (
	availableResolutions = []struct {
		w, h int
	}{
		{800, 600},
		{1280, 720},
		{1920, 1080},
		{2560, 1440},
	}
	currentResolution = 0
)

type Player struct {
	x, y       float64
	velX, velY float64
	isJumping  bool
	onSlope    bool
	color      color.Color
}

type Platform struct {
	x, y, width float64
	color       color.Color
	isSlope     bool
}

type Game struct {
	player         Player
	platforms      []Platform
	cameraX        float64
	gameOver       bool
	score          int
	lastPlatformX  float64
	startTimer     int
	isStarting     bool
	inMenu         bool
	inSettings     bool
	currentSetting int
	jumpKey        ebiten.Key
	screenWidth    int
	screenHeight   int
	currentSpeed   float64
}

func (g *Game) Update() error {
	if g.inMenu {
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			g.currentSetting = (g.currentSetting + 1) % 3
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.currentSetting = (g.currentSetting + 2) % 3
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if g.inSettings {
				switch g.currentSetting {
				case 0:
					currentResolution = (currentResolution + 1) % len(availableResolutions)
					g.screenWidth = availableResolutions[currentResolution].w
					g.screenHeight = availableResolutions[currentResolution].h
					ebiten.SetWindowSize(g.screenWidth, g.screenHeight)
				case 1:
					keys := inpututil.AppendPressedKeys(nil)
					if len(keys) > 0 {
						g.jumpKey = keys[0]
					}
				case 2:
					g.inSettings = false
				}
			} else {
				switch g.currentSetting {
				case 0:
					g.inMenu = false
					g.resetGame()
				case 1:
					g.inSettings = true
					g.currentSetting = 0
				case 2:
					g.inMenu = false
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			if g.inSettings {
				g.inSettings = false
			} else {
				g.inMenu = false
			}
		}
		return nil
	}

	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.resetGame()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.inMenu = true
		}
		return nil
	}

	if g.isStarting {
		g.startTimer--
		if g.startTimer <= 0 {
			g.isStarting = false
		}

		if inpututil.IsKeyJustPressed(g.jumpKey) && !g.player.isJumping {
			g.player.velY = jumpForce
			g.player.isJumping = true
			g.isStarting = false
		}
		return nil
	}

	// Apply speed increase every 30 points
	if g.score > 0 && g.score%speedIncreaseInterval == 0 {
		g.currentSpeed = baseSpeed * (1 + speedIncrease*float64(g.score/speedIncreaseInterval))
	}

	if inpututil.IsKeyJustPressed(g.jumpKey) && !g.player.isJumping {
		g.player.velY = jumpForce
		g.player.isJumping = true
	}

	g.player.velY += gravity
	g.player.x += g.currentSpeed

	// Slope platform handling
	g.player.onSlope = false
	for _, p := range g.platforms {
		if p.isSlope && g.player.x+playerSize > p.x && g.player.x < p.x+p.width {
			slopeRatio := (g.player.x - p.x) / p.width
			currentSlopeHeight := slopeRatio * slopeHeight
			platformTop := p.y - currentSlopeHeight

			if g.player.y+playerSize >= platformTop &&
				g.player.y+playerSize <= platformTop+20 {

				g.player.y = platformTop - playerSize
				g.player.velY = 0
				g.player.isJumping = false
				g.player.onSlope = true

				if slopeRatio < 0.9 {
					g.player.y += 0.8
				}

				if inpututil.IsKeyJustPressed(g.jumpKey) {
					g.player.velY = jumpForce
					g.player.isJumping = true
					g.player.onSlope = false
				}
				break
			}
		}
	}

	if !g.player.onSlope {
		g.player.y += g.player.velY
	}

	// Regular platform collisions
	g.player.isJumping = true
	for _, p := range g.platforms {
		if p.isSlope {
			continue
		}

		if g.player.y+playerSize >= p.y &&
			g.player.y+playerSize <= p.y+20 &&
			g.player.x+playerSize > p.x &&
			g.player.x < p.x+p.width &&
			g.player.velY >= 0 {

			g.player.y = p.y - playerSize
			g.player.velY = 0
			g.player.isJumping = false
			break
		}
	}

	if g.player.y > float64(g.screenHeight) {
		g.gameOver = true
	}

	// Platform generation
	if g.lastPlatformX < g.player.x+float64(g.screenWidth) {
		platformLevels := []struct {
			y     float64
			color color.Color
		}{
			{float64(g.screenHeight) - 50, color.RGBA{0, 255, 0, 255}},
			{float64(g.screenHeight) - 120, color.RGBA{0, 200, 255, 255}},
			{float64(g.screenHeight) - 190, color.RGBA{255, 255, 0, 255}},
			{float64(g.screenHeight) - 260, color.RGBA{255, 165, 0, 255}},
			{float64(g.screenHeight) - 330, color.RGBA{255, 0, 0, 255}},
		}

		currentLevel := 0
		minDist := math.MaxFloat64
		for i, level := range platformLevels {
			dist := math.Abs(g.player.y - (level.y - playerSize))
			if dist < minDist {
				minDist = dist
				currentLevel = i
			}
		}

		isSlope := false
		if (currentLevel >= 3) && rand.Intn(10) == 0 {
			isSlope = true
		}

		targetLevel := currentLevel
		if isSlope {
			targetLevel = 1 + rand.Intn(2)
		} else {
			if rand.Intn(2) == 0 && currentLevel > 0 {
				targetLevel = currentLevel - 1
			} else if currentLevel < len(platformLevels)-1 {
				targetLevel = currentLevel + 1
			}
		}

		newX := g.lastPlatformX + g.currentSpeed*50
		newWidth := 80.0 + rand.Float64()*70.0
		if isSlope {
			newWidth = slopeLength
		}

		g.platforms = append(g.platforms, Platform{
			x:       newX,
			y:       platformLevels[targetLevel].y,
			width:   newWidth,
			color:   platformLevels[targetLevel].color,
			isSlope: isSlope,
		})

		g.lastPlatformX = newX
	}

	// Remove passed platforms
	for i := 0; i < len(g.platforms); {
		if g.platforms[i].x+g.platforms[i].width < g.player.x-100 {
			g.platforms = append(g.platforms[:i], g.platforms[i+1:]...)
			g.score++
		} else {
			i++
		}
	}

	g.cameraX = g.player.x - 100
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 40, 255})

	if g.inMenu {
		title := "JUMP-JUMP-JUMP"
		options := []string{"START GAME", "SETTINGS", "EXIT TO GAME"}

		text.Draw(screen, title, basicfont.Face7x13, g.screenWidth/2-70, g.screenHeight/3, color.White)

		for i, opt := range options {
			var clr color.Color = color.White
			if i == g.currentSetting {
				clr = color.RGBA{R: 255, G: 255, B: 0, A: 255}
				text.Draw(screen, ">", basicfont.Face7x13, g.screenWidth/2-90, g.screenHeight/2+i*30, clr)
			}
			text.Draw(screen, opt, basicfont.Face7x13, g.screenWidth/2-50, g.screenHeight/2+i*30, clr)
		}

		if g.inSettings {
			settings := []string{
				fmt.Sprintf("RESOLUTION: %dx%d", g.screenWidth, g.screenHeight),
				fmt.Sprintf("JUMP KEY: %s", g.jumpKey.String()),
				"BACK",
			}

			for i, setting := range settings {
				var clr color.Color = color.White
				if i == g.currentSetting {
					clr = color.RGBA{R: 255, G: 255, B: 0, A: 255}
					text.Draw(screen, ">", basicfont.Face7x13, g.screenWidth/2-120, g.screenHeight/2+i*30, clr)
				}
				text.Draw(screen, setting, basicfont.Face7x13, g.screenWidth/2-100, g.screenHeight/2+i*30, clr)
			}
		}

		text.Draw(screen, "Use UP/DOWN and ENTER", basicfont.Face7x13, g.screenWidth/2-100, g.screenHeight-50, color.White)
		return
	}

	for _, p := range g.platforms {
		if p.isSlope {
			for i := 0; i < int(p.width); i++ {
				height := float64(i) / p.width * slopeHeight
				ebitenutil.DrawRect(screen, p.x-g.cameraX+float64(i), p.y-height, 1, 20+height, p.color)
			}
		} else {
			ebitenutil.DrawRect(screen, p.x-g.cameraX, p.y, p.width, 20, p.color)
		}
	}

	ebitenutil.DrawRect(screen, g.player.x-g.cameraX, g.player.y, playerSize, playerSize, g.player.color)

	if !g.isStarting {
		text.Draw(screen, fmt.Sprintf("SCORE: %d", g.score), basicfont.Face7x13, 20, 20, color.White)
		text.Draw(screen, fmt.Sprintf("SPEED: %.1f", g.currentSpeed), basicfont.Face7x13, 20, 40, color.White)
	} else {
		text.Draw(screen, "PRESS TO JUMP", basicfont.Face7x13, g.screenWidth/2-50, g.screenHeight/2, color.White)
		text.Draw(screen, fmt.Sprintf("(%s KEY)", g.jumpKey.String()), basicfont.Face7x13, g.screenWidth/2-40, g.screenHeight/2+20, color.White)
	}

	if g.gameOver {
		text.Draw(screen, "GAME OVER", basicfont.Face7x13, g.screenWidth/2-40, g.screenHeight/2-20, color.White)
		text.Draw(screen, fmt.Sprintf("SCORE: %d", g.score), basicfont.Face7x13, g.screenWidth/2-35, g.screenHeight/2+10, color.White)
		text.Draw(screen, "PRESS R TO RESTART", basicfont.Face7x13, g.screenWidth/2-70, g.screenHeight/2+40, color.White)
		text.Draw(screen, "ESC TO MENU", basicfont.Face7x13, g.screenWidth/2-40, g.screenHeight/2+70, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

func (g *Game) resetGame() {
	g.player = Player{
		x:     100,
		y:     float64(g.screenHeight) - 50 - playerSize,
		velX:  0,
		velY:  0,
		color: color.RGBA{255, 100, 100, 255},
	}

	g.platforms = []Platform{
		{x: 0, y: float64(g.screenHeight) - 50, width: 500, color: color.RGBA{0, 255, 0, 255}},
	}

	g.cameraX = 0
	g.gameOver = false
	g.score = 0
	g.lastPlatformX = 500
	g.startTimer = startPlatformDuration
	g.isStarting = true
	g.currentSpeed = baseSpeed
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		inMenu:       true,
		jumpKey:      ebiten.KeySpace,
		screenWidth:  availableResolutions[currentResolution].w,
		screenHeight: availableResolutions[currentResolution].h,
		currentSpeed: baseSpeed,
	}

	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("Geometry Dash-like Platformer")
	ebiten.SetWindowResizable(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
