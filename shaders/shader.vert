#version 330

uniform mat4 viewProj;
uniform mat4 model;

uniform float fogDist;

in vec4 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;
out float dist;

void main() {
  gl_Position = viewProj * model * vert;
  dist = max(min(length(gl_Position), fogDist), 0)/fogDist;
  fragTexCoord = vertTexCoord;
}