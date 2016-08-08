#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

uniform float fogDist;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;
out float dist;

void main() {
  gl_Position = projection * camera * model * vec4(vert, 1);
  dist = max(min(length(gl_Position), fogDist), 0)/fogDist;
  fragTexCoord = vertTexCoord;
}