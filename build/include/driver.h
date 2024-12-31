
// gcc -shared -fPIC -o test.so driver.c
// gcc -shared -fPIC -o test.dll driver.c

#ifndef DRIVER_H
#define DRIVER_H

#define DRV_SUCCESS    0


#include <stdio.h>
 #ifdef __cplusplus
extern "C" {
#endif

 int drv_start(char* param);
 int drv_stop(char* param);

 int drv_read_int(const char* param, int* value);
 int drv_read_double(const char* param, double* value);
 int drv_read_string(const char* param, char* data_buf, int buf_size);
 int drv_write_int(const char* param, int value);
 int drv_write_double(const char* param, double value);
 int drv_write_string(const char* param, char* data);

 #ifdef __cplusplus
}
#endif
#endif