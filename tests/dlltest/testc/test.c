

#include <stdio.h>
 #ifdef __cplusplus
extern "C" {
#endif
 int add(int a, int b) {
    return a + b;
}
 #ifdef __cplusplus
}
#endif
// ifndef MY_TEST_H
// define MY_TEST_H
// gcc -shared -fPIC -o test.so test.c
// endif