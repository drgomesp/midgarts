package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/logger"
)

type CharacterRenderSystem struct {
	log      *logger.Logger
	renderer *renderer.Renderer
	entities []Character
	scene    *core.Node
	camera   camera.ICamera
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
	s.scene.Add(e.GetRenderComponent().Graphic)
}

func (s *CharacterRenderSystem) AddByInterface(o ecs.Identifier) {
	e, _ := o.(Character)
	s.Add(e)
}

func (s *CharacterRenderSystem) Update(dt float32) {
	var rendered bool
	var err error

	if err = s.renderer.Render(s.scene, s.camera); err != nil {
		s.log.Fatal("could not update render system: %v\n", err)
	}

	s.log.Info("CharacterRenderSystem::Update(%v) rendered=(%v)", dt, rendered)
}

func (s *CharacterRenderSystem) Remove(e ecs.BasicEntity) {
	panic("implement me")
}
