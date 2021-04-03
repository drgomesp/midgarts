package main

import (
	"log"

	"github.com/EngoEngine/engo/common"
	"github.com/g3n/engine/graphic"
)

const (
	DefaultCharacterAnimationRate = 0.08
)

type CharacterAnimationComponent struct {
	Animator         *Animator
	Animations       map[string]*common.Animation
	CharacterSprite  *CharacterSprite
	CurrentAnimation *common.Animation
	CurrentFrame     uint32
	Rate             float32
	Index            int
	Change           float32
	DefaultAnimation *common.Animation
}

func NewCharacterAnimationComponent(characterSprite *CharacterSprite, rate float32) *CharacterAnimationComponent {
	return &CharacterAnimationComponent{
		Animator:        NewAnimator(characterSprite.Spritesheet),
		Animations:      make(map[string]*common.Animation),
		CharacterSprite: characterSprite,
		Rate:            rate,
	}
}

// SelectAnimationByName sets the current animation. The name must be
// registered.
func (ac *CharacterAnimationComponent) SelectAnimationByName(name string) {
	ac.CurrentAnimation = ac.Animations[name]
	ac.Index = 0
}

// SelectAnimationByAction sets the current animation.
// An nil action value selects the default animation.
func (ac *CharacterAnimationComponent) SelectAnimationByAction(action *common.Animation) {
	ac.CurrentAnimation = action
	ac.Index = 0
}

// AddDefaultAnimation adds an animation which is used when no other animation is playing.
func (ac *CharacterAnimationComponent) AddDefaultAnimation(action *common.Animation) {
	ac.AddAnimation(action)
	ac.DefaultAnimation = action
}

// AddAnimation registers an animation under its name, making it available
// through SelectAnimationByName.
func (ac *CharacterAnimationComponent) AddAnimation(action *common.Animation) {
	ac.Animations[action.Name] = action
}

// AddAnimations registers all given animations.
func (ac *CharacterAnimationComponent) AddAnimations(actions []*common.Animation) {
	for _, action := range actions {
		ac.AddAnimation(action)
	}
}

// Cell returns the drawable for the current frame.
func (ac *CharacterAnimationComponent) Cell() *graphic.Sprite {
	if len(ac.CurrentAnimation.Frames) == 0 {
		log.Println("No frame data for this animation. Selecting zeroth drawable. If this is incorrect, add an action to the animation.")
		return ac.CharacterSprite.Spritesheet.SpriteAt(0)
	}
	idx := uint32(ac.CurrentAnimation.Frames[ac.Index])
	ac.CurrentFrame = idx
	ac.Animator.CurrentFrame = idx
	return ac.CharacterSprite.Spritesheet.SpriteAt(idx)
}

// NextFrame advances the current animation by one frame.
func (ac *CharacterAnimationComponent) NextFrame() {
	if len(ac.CurrentAnimation.Frames) == 0 {
		log.Println("No frame data for this animation")
		return
	}

	ac.Index++
	ac.Change = 0
	if ac.Index >= len(ac.CurrentAnimation.Frames) {
		ac.Index = 0

		if !ac.CurrentAnimation.Loop {
			ac.CurrentAnimation = nil
			return
		}
	}
}
