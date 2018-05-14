
#include<stdio.h>
#include<stdlib.h>

int main() {
  printf("hello world\n");
  printf("uid %d\n", getuid());
  sleep(1);
  int i = 0;
  while(1) {
    i++;
  }
  printf("%d\n", i);
  return 0;
}
