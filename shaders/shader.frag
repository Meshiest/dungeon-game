#version 330 core
uniform sampler2D tex;

in vec2 fragTexCoord;
in float dist;

out vec4 color;
void main() {
  color = mix(texture(tex, fragTexCoord), vec4(0, 0, 0, 1), dist);
}