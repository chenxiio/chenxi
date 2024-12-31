#include<stdio.h>
#include<dlfcn.h>

int main()
{
	void* handle;
	int (*addfunc)(int,int);
	char* error;
	
	/* handle = dlopen("/lib/libm-2.6.1.so", RTLD_NOW); linux */
	handle = dlopen("./test.so", RTLD_LAZY);
	if(handle == NULL){
		printf("open lib error: %s\n", dlerror());
		return -1;
	}
	
	addfunc = dlsym(handle, "add");
	if(NULL != (error = dlerror())){
		printf("symbol sin not found , error: %s\n",error);
		return -1;
	}
	
	printf("%d\n",addfunc(1 ,2));
	dlclose(handle);
	return 0;
}

// root@ubuntu-admin-a1:/home/6Chapter# gcc -o dyTest dyTest.c -ldl
// root@ubuntu-admin-a1:/home/6Chapter# 
// root@ubuntu-admin-a1:/home/6Chapter# ./dyTest 
// 1.000000
// root@ubuntu-admin-a1:/home/6Chapter# 
