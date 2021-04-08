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
uniform mat4 mvp;

out vec3 fragColor;

void main() {
	gl_Position = mvp * vec4(position, 1.0);
	fragColor = vertexColor;
}` + "\x00"

	fragmentShaderSource = `
#version 330 core
out vec4 outColor;
in vec3 fragColor;
void main() {
	outColor = vec4(fragColor, 1.0);
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
