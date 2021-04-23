package main

import (
	"math"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/logger"
	"github.com/g3n/engine/window"
	"github.com/project-midgard/midgarts/internal/entity"
	"github.com/project-midgard/midgarts/internal/system"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
)

type CharacterRenderSystem struct {
	log      *logger.Logger
	renderer *renderer.Renderer
	entities []Character

	scene  *core.Node
	camera camera.ICamera

	spritesOnScene []int
}

func NewCharacterRenderSystem(log *logger.Logger, renderer *renderer.Renderer, scene *core.Node, camera camera.ICamera) *CharacterRenderSystem {
	return &CharacterRenderSystem{
		log:      log,
		renderer: renderer,
		scene:    scene,
		camera:   camera,
		entities: make([]Character, 0),
	}
}

func (s *CharacterRenderSystem) Add(e Character) {
	s.entities = append(s.entities, e)
	s.scene.AddAt(int(e.GetBasicEntity().ID()), e.GetCharacterRenderComponent().CharacterSprite.Spritesheet.SpriteAt(0))
}

func (s *CharacterRenderSystem) AddByInterface(o ecs.Identifier) {
	e, _ := o.(Character)
	s.Add(e)
}

func (s *CharacterRenderSystem) Update(dt float32) {
	var err error

	for _, e := range s.entities {
		if KeyState.Pressed(window.KeyUp) && KeyState.Pressed(window.KeyRight) {
			e.SetState(statetype.Walking)
			e.SetDirection(directiontype.NorthEast)
		} else if KeyState.Pressed(window.KeyUp) && KeyState.Pressed(window.KeyLeft) {
			e.SetState(statetype.Walking)
			e.SetDirection(directiontype.NorthWest)
		} else if KeyState.Pressed(window.KeyDown) && KeyState.Pressed(window.KeyRight) {
			e.SetState(statetype.Walking)
			e.SetDirection(directiontype.SouthEast)
		} else if KeyState.Pressed(window.KeyDown) && KeyState.Pressed(window.KeyLeft) {
			e.SetState(statetype.Walking)
			e.SetDirection(directiontype.SouthWest)
		} else if KeyState.Pressed(window.KeyUp) {
			e.SetState(statetype.Walking)
			e.SetDirection(directiontype.North)
		} else if KeyState.Pressed(window.KeyRight) {
			e.SetState(statetype.Walking)
			e.SetDirection(directiontype.East)
		} else if KeyState.Pressed(window.KeyDown) {
			e.SetState(statetype.Walking)
			e.SetDirection(directiontype.South)
		} else if KeyState.Pressed(window.KeyLeft) {
			e.SetState(statetype.Walking)
			e.SetDirection(directiontype.West)
		} else {
			e.SetState(statetype.Idle)
		}

		if e.GetCharacterAnimationComponent().CurrentAnimation == nil {
			if e.GetCharacterAnimationComponent().DefaultAnimation == nil {
				continue
			}

			e.GetCharacterAnimationComponent().
				SelectAnimationByAction(e.GetCharacterAnimationComponent().DefaultAnimation)
		}

		// TODO this only works for character's direction, not action
		actionIndex := actionindex.GetActionIndex(e.GetState())
		idx := int(actionIndex) + (int(e.GetDirection())+directiontype.DirectionTable[system.FixedCameraDirection])%8
		action := e.GetCharacterRenderComponent().CharacterSprite.ActionFile.Actions[idx]

		frameCount := len(action.Frames)
		timeNeededForOneFrame := int64(action.Delay.Seconds() * (1.0 / system.FPSMultiplier))
		timeNeededForOneFrame = int64(math.Max(float64(timeNeededForOneFrame), 100))

		elapsedTime := time.Since(e.GetCurrentAction().AnimationStartedAt).Milliseconds() - int64(dt)
		realIndex := elapsedTime / timeNeededForOneFrame

		var frameIndex int64
		switch e.GetPlayMode() {
		case actionplaymode.Repeat:
			frameIndex = realIndex % int64(frameCount)
			break
		}

		if frame := action.Frames[frameIndex]; len(frame.Layers) == 0 {
			continue
		}

		var frames []int
		for _, f := range action.Frames {
			frames = append(frames, int(f.Layers[0].SpriteFrameIndex))
		}

		e.SetCurrentAction(entity.NewCharacterAction(actionindex.GetActionIndex(e.GetState())))
		e.GetCurrentAction().SetFrames(frames)
		anim := &common.Animation{Name: e.GetCurrentAction().Name, Frames: e.GetCurrentAction().Frames}
		e.GetCharacterAnimationComponent().AddAnimations([]*common.Animation{anim})
		e.GetCharacterAnimationComponent().AddDefaultAnimation(anim)
		e.GetCharacterAnimationComponent().CurrentAnimation = anim
		e.GetCharacterAnimationComponent().Change += dt

		if e.GetCharacterAnimationComponent().Change >= e.GetCharacterAnimationComponent().Rate {
			e.GetCharacterAnimationComponent().CurrentFrame = uint32(idx * int(frameIndex))

			if e.GetCharacterAnimationComponent().Index >= len(e.GetCurrentAction().Frames) {
				e.GetCharacterAnimationComponent().Index = 0
			}

			e.GetCharacterAnimationComponent().Cell()
			e.GetCharacterAnimationComponent().NextFrame()
		}

		e.GetCharacterAnimationComponent().Animator.Update(time.Now())
	}

	if err = s.renderer.Render(s.scene, s.camera); err != nil {
		s.log.Fatal("could not update render system: %v\n", err)
	}
}

func (s *CharacterRenderSystem) Remove(e ecs.BasicEntity) {
	s.scene.RemoveAt(int(e.GetBasicEntity().ID()))
}
