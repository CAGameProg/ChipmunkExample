package main

import (
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	sf "github.com/zyedidia/sfml/v2.3/sfml"
)

const (
	screenWidth  = 800
	screenHeight = 600

	ballRadius = 30
	ballMass   = 1
)

type Player struct {
	*chipmunk.Body
	*sf.Sprite

	keys []sf.KeyCode
}

func NewPlayer(pos sf.Vector2f, keys []sf.KeyCode, texture *sf.Texture) *Player {
	p := new(Player)
	sprite := sf.NewSprite(texture)
	size := sprite.GetGlobalBounds()
	sprite.SetOrigin(sf.Vector2f{size.Width / 2, size.Height / 2})
	sprite.SetPosition(pos)
	p.Sprite = sprite
	p.keys = keys

	rect := chipmunk.NewBox(vect.Vector_Zero, vect.Float(size.Width), vect.Float(size.Height))
	// rect.SetElasticity(0.95)

	body := chipmunk.NewBody(vect.Float(ballMass), rect.Moment(float32(ballMass)))
	body.SetPosition(vect.Vect{vect.Float(pos.X), vect.Float(pos.Y)})
	// body.SetAngle(vect.Float(rand.Float32() * 2 * math.Pi))

	body.AddShape(rect)
	space.AddBody(body)
	p.Body = body

	return p
}

func (p *Player) Update() {
	// b.Sprite.SetRotation(0)
	p.Body.SetAngle(0)
	// p.Shape.Get
	// p.SetPosition(body)
	pos := p.Body.Position()
	p.Sprite.SetPosition(sf.Vector2f{float32(pos.X), float32(-pos.Y)})
	v := p.Velocity()

	if sf.KeyboardIsKeyPressed(p.keys[0]) {
		p.SetVelocity(float32(v.X), 400)
	}
	if sf.KeyboardIsKeyPressed(p.keys[2]) {
		p.AddVelocity(-10, 0)
	}
	if sf.KeyboardIsKeyPressed(p.keys[3]) {
		p.AddVelocity(10, 0)
	}
}

type Ball struct {
	*chipmunk.Body
	*sf.Sprite
}

func NewBall(pos sf.Vector2f, texture *sf.Texture) *Ball {
	b := new(Ball)
	sprite := sf.NewSprite(texture)
	size := sprite.GetGlobalBounds()
	sprite.SetOrigin(sf.Vector2f{size.Width / 2, size.Height / 2})
	sprite.SetPosition(pos)
	b.Sprite = sprite

	ball := chipmunk.NewCircle(vect.Vector_Zero, float32(size.Width/2))
	ball.SetElasticity(0.95)

	body := chipmunk.NewBody(vect.Float(ballMass), ball.Moment(float32(ballMass)))
	body.SetPosition(vect.Vect{vect.Float(pos.X), vect.Float(-pos.Y)})
	// body.SetAngle(vect.Float(rand.Float32() * 2 * math.Pi))

	body.AddShape(ball)
	space.AddBody(body)
	b.Body = body
	// balls = append(balls, b)

	return b
}

func (b *Ball) Update() {
	pos := b.Body.Position()
	b.Sprite.SetPosition(sf.Vector2f{float32(pos.X), float32(-pos.Y)})
	angle := b.Body.Angle()
	b.Sprite.SetRotation(180.0 / math.Pi * float32(-angle))
}

var (
	space *chipmunk.Space

	player      *Player
	balls       []*Ball
	staticLines []*chipmunk.Shape
)

func step() {
	space.Step(vect.Float(1.0 / 60.0))
}

// createBodies sets up the chipmunk space and static bodies
func createWorld() {
	space = chipmunk.NewSpace()
	space.Gravity = vect.Vect{0, -900}

	staticBody := chipmunk.NewBodyStatic()
	staticLines = []*chipmunk.Shape{
		chipmunk.NewSegment(vect.Vect{0, -600}, vect.Vect{800.0, -600}, 0),
		chipmunk.NewSegment(vect.Vect{0, -600}, vect.Vect{0, 0}, 0),
		chipmunk.NewSegment(vect.Vect{800, -600}, vect.Vect{800.0, 0}, 0),
	}
	for _, segment := range staticLines {
		// segment.SetElasticity(0.6)
		staticBody.AddShape(segment)
	}
	space.AddBody(staticBody)
}

func main() {
	runtime.LockOSThread()

	window := sf.NewRenderWindow(sf.VideoMode{screenWidth, screenHeight, 32}, "Space Shooter", sf.StyleDefault, nil)
	window.SetFramerateLimit(60)

	createWorld()

	player := NewPlayer(sf.Vector2f{400, 0}, []sf.KeyCode{sf.KeyUp, sf.KeyDown, sf.KeyLeft, sf.KeyRight}, sf.NewTexture("mario.png"))

	lock := sync.RWMutex{}

	go func() {
		for {
			lock.Lock()
			balls = append(balls, NewBall(sf.Vector2f{400, 0}, sf.NewTexture("smiley.png")))
			lock.Unlock()
			time.Sleep(500 * time.Millisecond)
		}
	}()

	for window.IsOpen() {
		if event := window.PollEvent(); event != nil {
			switch event.Type {
			case sf.EventClosed:
				window.Close()
			}
		}

		step()
		player.Update()
		lock.RLock()
		for _, ball := range balls {
			ball.Update()
		}
		lock.RUnlock()

		window.Clear(sf.ColorWhite)

		window.Draw(player.Sprite)
		for _, ball := range balls {
			window.Draw(ball.Sprite)
		}

		window.Display()
	}
}
