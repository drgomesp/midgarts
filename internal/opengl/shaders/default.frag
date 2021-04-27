#version 330 core

in vec3 fragColor;
in vec2 texCoords;

uniform sampler2D tex;

void main() {
    vec2 var_TexCoords = texCoords;

    //gl_FragColor = vec4(fragColor, 1.0);
    //gl_FragColor = texture(tex, var_TexCoords);
    gl_FragColor = texture(tex, var_TexCoords)  * vec4(fragColor, 1.0);
}
