package opengl

const (
	vertexShaderSource = `
#version 330 core

layout(location = 0) in vec3 VertexPosition;
layout(location = 1) in vec3 VertexColor;
layout(location = 2) in vec2 VertexTexCoord;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform mat4 rotation;
uniform vec2 size;
uniform vec2 offset;

out vec3 fragColor;
out vec2 texCoords;

void main() {
	vec4 pos = vec4(VertexPosition.x * size.x, VertexPosition.y * size.y, 0.0, 1.0);
	pos.x += offset.x;
	pos.y -= offset.y;

	mat4 modelView = view * model;

	modelView[0].xyz = vec3( 1.0, 0.0, 0.0 );
    modelView[1].xyz = vec3( 0.0, 1.0, 0.0 );
    modelView[2].xyz = vec3( 0.0, 0.0, 1.0 );

	gl_Position = projection * modelView * (rotation * pos);

	fragColor = VertexColor;
	texCoords = VertexTexCoord;
}` + "\x00"

	fragmentShaderSource = `
#version 330 core

in vec3 fragColor;
in vec2 texCoords;

uniform sampler2D tex;

void main() {
	vec2 var_TexCoords = texCoords;

	//gl_FragColor = vec4(fragColor, 1.0);
	//gl_FragColor = texture(tex, var_TexCoords);
	gl_FragColor = texture(tex, var_TexCoords)  * vec4(fragColor, 1.0);
}` + "\x00"
)
