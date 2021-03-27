package system

import (
	"log"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/engo/math"
	"github.com/project-midgard/midgarts/pkg/client"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
)

const FPSMultiplier = 1.0
const FixedCameraDirection = 6

var OneSpritePixelSizeIn3D float32 = 1.0 / 35.0

var DirectionTable = [8]int{6, 5, 4, 3, 2, 1, 0, 7}

type CharacterAnimationSystem struct {
	characters map[string]*client.CharacterEntity
}

func NewCharacterAnimationSystem() *CharacterAnimationSystem {
	return &CharacterAnimationSystem{map[string]*client.CharacterEntity{}}
}

func (s *CharacterAnimationSystem) Add(char *client.CharacterEntity) {
	s.characters[char.UUID()] = char
}

func (s *CharacterAnimationSystem) AddByInterface(i ecs.Identifier) {
	o, _ := i.(*client.CharacterEntity)
	s.Add(o)
}

func (s *CharacterAnimationSystem) Update(dt float32) {
	for _, char := range s.characters {
		if engo.Input.Button("Top").Down() && engo.Input.Button("Right").Down() {
			char.State = statetype.Walking
			char.Direction = directiontype.NorthEast
		} else if engo.Input.Button("Top").Down() && engo.Input.Button("Left").Down() {
			char.State = statetype.Walking
			char.Direction = directiontype.NorthWest
		} else if engo.Input.Button("Bot").Down() && engo.Input.Button("Right").Down() {
			char.State = statetype.Walking
			char.Direction = directiontype.SouthEast
		} else if engo.Input.Button("Bot").Down() && engo.Input.Button("Left").Down() {
			char.State = statetype.Walking
			char.Direction = directiontype.SouthWest
		} else if engo.Input.Button("Top").Down() {
			char.State = statetype.Walking
			char.Direction = directiontype.North
		} else if engo.Input.Button("Right").Down() {
			char.State = statetype.Walking
			char.Direction = directiontype.East
		} else if engo.Input.Button("Bot").Down() {
			char.State = statetype.Walking
			char.Direction = directiontype.South
		} else if engo.Input.Button("Left").Down() {
			char.State = statetype.Walking
			char.Direction = directiontype.West
		} else {
			//char.Direction = directiontype.South
			char.State = statetype.Idle
		}

		if char.CharacterAnimationComponent.CurrentAnimation == nil {
			if char.CharacterAnimationComponent.DefaultAnimation == nil {
				continue
			}

			char.
				CharacterAnimationComponent.
				SelectAnimationByAction(char.CharacterAnimationComponent.DefaultAnimation)
		}

		actionIndex := actionindex.GetActionIndex(char.State)
		idx := int(actionIndex) + (int(char.Direction)+DirectionTable[FixedCameraDirection])%8
		action := char.ActionFile.Actions[idx]

		frameCount := len(action.Frames)
		timeNeededForOneFrame := int64(action.Delay.Seconds() * (1.0 / FPSMultiplier))
		timeNeededForOneFrame = int64(math.Max(float32(timeNeededForOneFrame), 100))

		elapsedTime := time.Since(char.CurrentAction.AnimationStartedAt).Milliseconds() - int64(dt)
		realIndex := elapsedTime / timeNeededForOneFrame

		var frameIndex int64
		switch char.PlayMode {
		case actionplaymode.Repeat:
			frameIndex = realIndex % int64(frameCount)
			break
		}

		frame := action.Frames[frameIndex]
		layer := frame.Layers[0]

		if len(frame.Layers) == 0 {
			continue
		}

		isMain := true
		providedOffset := [2]int32{0, 0}
		var offset [2]int32

		if !isMain {
			if len(frame.Positions) > 0 {
				offset = [2]int32{
					providedOffset[0] - frame.Positions[frameIndex][0],
					providedOffset[1] - frame.Positions[frameIndex][1],
				}
			}
		} else {
			offset = [2]int32{0, 0}
		}

		offset[0] = layer.Position[0] + offset[0]
		offset[1] = layer.Position[1] + offset[1]

		var frames []int
		for _, f := range action.Frames {
			frames = append(frames, int(f.Layers[0].SpriteFrameIndex))
		}

		char.CurrentAction = client.NewCharacterAction(actionindex.GetActionIndex(char.State))
		char.CurrentAction.SetFrames(frames)
		anim := &common.Animation{Name: char.CurrentAction.Name, Frames: char.CurrentAction.Frames}
		char.CharacterAnimationComponent.AddAnimations([]*common.Animation{anim})
		char.CharacterAnimationComponent.AddDefaultAnimation(anim)
		char.CharacterAnimationComponent.CurrentAnimation = anim
		char.CharacterAnimationComponent.Change += dt

		if char.CharacterAnimationComponent.Change >= char.CharacterAnimationComponent.Rate {
			//char.CharacterAnimationComponent.CurrentFrame = int(frameIndex)
			char.CharacterAnimationComponent.CurrentFrame = idx * int(frameIndex)

			if char.CharacterAnimationComponent.Index >= len(char.CurrentAnimation.Frames) {
				char.CharacterAnimationComponent.Index = 0
			}

			var posOffset engo.Point

			posOffset = engo.Point{
				X: float32(offset[0]) * OneSpritePixelSizeIn3D,
				Y: float32(offset[1]) * OneSpritePixelSizeIn3D,
			}

			w := char.SpritesheetResource.Spritesheet.Cell(0).Width()

			if layer.IsMirror {
				if char.Scale.X > 0 {
					char.Scale.Set(-1, 1)
					char.Position.Set(char.Position.X+w, char.Position.Y)
				}
			} else {
				if char.Scale.X < 0 {
					char.Scale.Set(1, 1)
					char.Position.Set(char.Position.X-w, char.Position.Y)
				}
			}

			//spriteShader := char.RenderComponent.Shader().(*common.SpriteShader)
			//spriteShader.SetModelMatrix((&engo.Matrix{}).Translate(char.Position.X, char.Position.Y))
			//spriteShader.SetSpriteSize(engo.Point{
			//	X: float32(flippedWidth),
			//	Y: float32(layer.Height),
			//})
			//char.Rotation = (&engo.Matrix{Val: [9]float32{
			//	1.0, 0.0, 0.0,
			//	0.0, 1.0, 0.0,
			//	0.0, 0.0, 1.0,
			//}}).RotationComponent()

			//char.Position.MultiplyMatrixVector(mat.Translate(-0.5*char.Scale.X, -0.5*char.Scale.Y))

			log.Printf("pos=%v offset=%v\n", char.Position, posOffset)

			char.RenderComponent.Drawable = char.CharacterAnimationComponent.Cell()
			char.CharacterAnimationComponent.NextFrame()
		}
	}
}

func (s CharacterAnimationSystem) Remove(e ecs.BasicEntity) {
	panic("implement me")
}
