package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 1500
	screenHeight = 900
	paddleWidth  = 10
	paddleHeight = 100
	ballRadius   = 10

	paddleSpeed = 7 // vitesse de déplacement des raquettes
)

type Game struct {
	y1, y2                  float32
	xball, yball, ballSpeed float32
	vx, vy                  float32
	lose                    int
}

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	// initialiser le random
	rand.Seed(time.Now().UnixNano())
}

func (g *Game) Update() error {
	g.vx = g.ballSpeed
	g.vx = g.ballSpeed / 2
	if g.lose == 0 {
		// Déplacement des raquettes
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			g.y1 -= paddleSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.y1 += paddleSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			g.y2 -= paddleSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.y2 += paddleSpeed
		}
	}

	// Déplacement de la balle
	if g.vx < 0 {
		g.xball -= (g.vy + g.ballSpeed) - 2
		g.yball -= (g.vy + (g.ballSpeed / 2)) - 2
	}
	if g.vx >= 0 {
		g.xball += (g.vy + g.ballSpeed) - 2
		g.yball += (g.vy + (g.ballSpeed / 2)) - 2
	}

	// Rebonds haut/bas avec rayon
	if g.yball-ballRadius <= 0 || g.yball+ballRadius >= screenHeight {
		g.vy = -g.vy
		g.vx = -g.vx
	}

	// Collision raquette gauche
	if g.vx < 0 && g.xball-ballRadius <= 20 &&
		g.yball >= g.y1 && g.yball <= g.y1+paddleHeight {

		g.vx = -g.vx

		// on calcule où la balle a touché la raquette
		hitPos := (g.yball - g.y1) - paddleHeight/2
		g.vy += hitPos * 0.05
		if g.ballSpeed == 10 {
			g.ballSpeed = 6
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.ballSpeed = 10
		}
	}

	// Collision raquette droite
	if g.vx > 0 && g.xball+ballRadius >= 1420 &&
		g.yball >= g.y2 && g.yball <= g.y2+paddleHeight {

		g.vx = -g.vx

		hitPos := (g.yball - g.y2) - paddleHeight/2
		g.vy += hitPos * 0.05
		if g.ballSpeed == 10 {
			g.ballSpeed = 6
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.ballSpeed = 10
		}
	}

	// Détection victoire/défaite
	if g.lose == 0 {
		if g.xball < 0 {
			g.lose = 1
		}
		if g.xball > screenWidth {
			g.lose = 2
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Debug
	ebitenutil.DebugPrint(screen, "Pong !")

	// Raquettes
	green := color.RGBA{0, 255, 0, 255}
	vector.DrawFilledRect(screen, 10, g.y1, paddleWidth, paddleHeight, green, true)
	vector.DrawFilledRect(screen, 1420, g.y2, paddleWidth, paddleHeight, green, true)

	// Balle
	white := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledCircle(screen, g.xball, g.yball, ballRadius, white, true)

	// Dead screens
	if g.lose == 1 {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(1500/4), float64(900/2))
		op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
		text.Draw(screen, "PLAYER 2 WIN", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   75,
		}, op)
	}
	if g.lose == 2 {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(1500/4), float64(900/2))
		op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
		text.Draw(screen, "PLAYER 1 WIN", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   75,
		}, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong")

	// direction aléatoire pour la balle
	startVy := float32(rand.Intn(7) - 3) // entre -3 et +3
	if startVy == 0 {
		startVy = 2 // éviter 0 (balle trop plate)
	}

	g := Game{
		y1:        350,
		y2:        350,
		xball:     screenWidth / 2,
		yball:     screenHeight / 2,
		ballSpeed: 6,
		vx:        6,
		vy:        startVy,
	}

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
