package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	vertexShaderSource = `
#version 330 core

layout(location = 0) in vec3 position;
layout(location = 1) in vec3 vertexColor;
layout(location = 2) in vec2 vertexCoords;

uniform mat4 mvp;

out vec3 fragColor;
out vec2 texCoords;

void main() {
	gl_Position = mvp * vec4(position, 1.0);
	fragColor = vertexColor;
	texCoords = vertexCoords;
}` + "\x00"

	fragmentShaderSource = `
#version 330 core

in vec3 fragColor;
in vec2 texCoords;

uniform sampler2D diffuse;

void main() {
	gl_FragColor = texture(diffuse, texCoords);
}` + "\x00"
)

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
