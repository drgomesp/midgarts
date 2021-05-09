#version 330 core

in vec3 fragColor;
in vec2 texCoords;

uniform sampler2D tex;

void main() {
    vec2 var_TexCoords = texCoords;
    vec4 texColor = texture(tex, var_TexCoords);

    if(texColor.a < 0.1)
        discard;

    gl_FragColor = texColor * vec4(fragColor, 1.0);
}
