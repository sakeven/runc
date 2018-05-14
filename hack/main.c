#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/mount.h>
#include <stdio.h>
#include <sched.h>
#include <signal.h>
#include <unistd.h>
#include <sys/resource.h>
 
int main() {
	system("ls -l proc");
	system("cat /proc/1/status");
    struct rlimit LIM; // time limit, file limit& memory limit
    // time limit
    
    LIM.rlim_cur =  1;
    LIM.rlim_max = LIM.rlim_cur;

    setrlimit(RLIMIT_CPU, &LIM);
	execl("./main", "./main",  (char *)NULL);
	return 0;
}
