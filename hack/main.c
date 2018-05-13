#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/mount.h>
#include <stdio.h>
#include <sched.h>
#include <signal.h>
#include <unistd.h>
 
int main() {
	system("ls -l proc");
	system("cat /proc/1/status");
	system("cat /proc/4/status");
	system("cat /proc/5/status");
	//execl("./test.sh", "./test.sh",  (char *)NULL);
	return 0;
}
