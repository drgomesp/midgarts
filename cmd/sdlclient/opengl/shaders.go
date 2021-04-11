package opengl

const (
	vertexShaderSource = `
#version 330 core

layout(location = 0) in vec3 VertexPosition;
layout(location = 1) in vec3 VertexColor;
layout(location = 2) in vec2 VertexTexCoord;

uniform mat4 mvp;

out vec3 fragColor;
out vec2 texCoords;

void main() {
	gl_Position = mvp * vec4(VertexPosition, 1.0);
	fragColor = VertexColor;
	texCoords = VertexTexCoord;
}` + "\x00"

	fragmentShaderSource = `
#version 330 core

in vec3 fragColor;
in vec2 texCoords;

uniform sampler2D diffuse;

void main() {
	//gl_FragColor = vec4(fragColor, 1.0);
	//gl_FragColor = texture(diffuse, texCoords);

	vec2 var_TexCoords = texCoords;
	gl_FragColor = texture(diffuse, var_TexCoords) * vec4(fragColor, 1.0);
}` + "\x00"
)
