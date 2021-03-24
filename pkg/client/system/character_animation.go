package system

import (
	"time"

	"github.com/EngoEngine/engo/math"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/pkg/client"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
)

const FPSMultiplier = 1.0
const FixedCameraDirection = 6

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
		if engo.Input.Button("Top").JustPressed() {
			char.Direction = directiontype.North
		} else if engo.Input.Button("Right").JustPressed() {
			char.Direction = directiontype.East
		} else if engo.Input.Button("Bot").JustPressed() {
			char.Direction = directiontype.South
		} else if engo.Input.Button("Left").JustPressed() {
			char.Direction = directiontype.West
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

		elapsedTime := int64(time.Since(char.CurrentAction.AnimationStartedAt).Milliseconds()) - int64(dt)
		realIndex := elapsedTime / timeNeededForOneFrame

		var frameIndex int64
		switch char.PlayMode {
		case actionplaymode.Repeat:
			frameIndex = realIndex % int64(frameCount)
			break
		}

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
			char.CharacterAnimationComponent.CurrentFrame = char.CurrentAnimation.Frames[frameIndex]
			char.RenderComponent.Drawable = char.CharacterAnimationComponent.Cell()
			char.CharacterAnimationComponent.NextFrame()
		}
	}
}

func (s CharacterAnimationSystem) Remove(e ecs.BasicEntity) {
	panic("implement me")
}
