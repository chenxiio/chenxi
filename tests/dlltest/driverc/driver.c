#include "../../../build/include/driver.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>


 int add(int a, int b) {
    return a + b;
}



 int drv_start(char* param) {
    // 实现drv_start接口
    return rand();
}
 int drv_stop(char* param) {
    // 实现drv_stop接口
    return rand();
}
 int drv_read_int(const char* param, int* value) {
    // 实现drv_read_int接口
    *value = rand();
    return rand();
}
 int drv_read_double(const char* param, double* value) {
    // 实现drv_read_double接口
    *value = (double)rand() / RAND_MAX;
    return rand();
}
 int drv_read_string(const char* param, char* data_buf, int buf_size) {
    // 实现drv_read_string接口
    char* str = "Hello World";
    int len = strlen(str);
    if (len > buf_size) {
        return -1;
    }
    memcpy(data_buf, str, len);
    return rand();
}
 int drv_write_int(const char* param, int value) {
    // 实现drv_write_int接口
    return rand();
}
 int drv_write_double(const char* param, double value) {
    // 实现drv_write_double接口
    return rand();
}
 int drv_write_string(const char* param, char* data) {
    // 实现drv_write_string接口
    return rand();
}

//  int main() {
//     srand(time(NULL));
//     // 调用接口，并随机输出返回值
//     printf("drv_start: %d\n", drv_start("param"));
//     printf("drv_stop: %d\n", drv_stop("param"));
//     int int_value;
//     printf("drv_read_int: %d, %d\n", drv_read_int("param", &int_value), int_value);
//     double double_value;
//     printf("drv_read_double: %d, %f\n", drv_read_double("param", &double_value), double_value);
//     char str_buf[20];
//     printf("drv_read_string: %d, %s\n", drv_read_string("param", str_buf, 20), str_buf);
//     printf("drv_write_int: %d\n", drv_write_int("param", rand()));
//     printf("drv_write_double: %d\n", drv_write_double("param", (double)rand() / RAND_MAX));
//     printf("drv_write_string: %d\n", drv_write_string("param", "Hello"));
//     return 0;
// }