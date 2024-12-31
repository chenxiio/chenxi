
// gcc -shared -fPIC -o ../../../bin/drivers/testdriver.so driver.c
// gcc -shared -fPIC -o ../../../bin/drivers/testdriver.dll driver.c
#include "../../../../../build/include/driver.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>



 typedef struct {
    char* key;
    void* value;
} KeyValue;
 typedef struct {
    KeyValue* kvArray;
    size_t size;
} KeyValueArray;
 KeyValueArray* kvArr;
 void initKeyValueArray(KeyValueArray* kvArr) {
    kvArr->kvArray = NULL;
    kvArr->size = 0;
}
 void addKeyValue(KeyValueArray* kvArr, const char* key, void* value) {
    kvArr->kvArray = realloc(kvArr->kvArray, (kvArr->size + 1) * sizeof(KeyValue));
    kvArr->kvArray[kvArr->size].key = strdup(key);
    kvArr->kvArray[kvArr->size].value = value;
    kvArr->size++;
}
 void freeKeyValueArray(KeyValueArray* kvArr) {
    for (size_t i = 0; i < kvArr->size; i++) {
        free(kvArr->kvArray[i].key);
    }
    free(kvArr->kvArray);
    initKeyValueArray(kvArr);
}


 int add(int a, int b) {
    return a + b;
}

 int drv_start(char* param) {
    // 实现drv_start接口
        kvArr = malloc(sizeof(KeyValueArray));
     initKeyValueArray(kvArr);
    return 0;
}
 int drv_stop(char* param) {
    // 实现drv_stop接口
    freeKeyValueArray(kvArr);
    free(kvArr);
    return 0;
}
 int drv_read_int(const char* param, int* value) {
    // 实现drv_read_int接口
    for (size_t i = 0; i < kvArr->size; i++) {
        if (strcmp(kvArr->kvArray[i].key, param) == 0) {
            *value = *(int*)kvArr->kvArray[i].value;
            return 0;
        }
    }
    *value = rand();
    return 0;
}
 int drv_read_double(const char* param, double* value) {
    // 实现drv_read_double接口
    for (size_t i = 0; i < kvArr->size; i++) {
        if (strcmp(kvArr->kvArray[i].key, param) == 0) {
            *value = *(double*)kvArr->kvArray[i].value;
            return 0;
        }
    }
    *value = (double)rand() / RAND_MAX;
    return 0;
}
 int drv_read_string(const char* param, char* data_buf, int buf_size) {
    // 实现drv_read_string接口
    for (size_t i = 0; i < kvArr->size; i++) {
        if (strcmp(kvArr->kvArray[i].key, param) == 0) {
            char* str = (char*)kvArr->kvArray[i].value;
            int len = strlen(str);
            if (len > buf_size) {
                return -1;
            }
            memcpy(data_buf, str, len);
            return 0;
        }
    }
    char* str = "Hello World";
    int len = strlen(str);
    if (len > buf_size) {
        return -1;
    }
    memcpy(data_buf, str, len);
    return 0;
}
 int drv_write_int(const char* param, int value) {
    // 实现drv_write_int接口
    int* pValue = malloc(sizeof(int));
    *pValue = value;
    addKeyValue(kvArr, param, pValue);
    return 0;
}
 int drv_write_double(const char* param, double value) {
    // 实现drv_write_double接口
    double* pValue = malloc(sizeof(double));
    *pValue = value;
    addKeyValue(kvArr, param, pValue);
    return 0;
}
 int drv_write_string(const char* param, char* data) {
    // 实现drv_write_string接口
    char* pValue = strdup(data);
    addKeyValue(kvArr, param, pValue);
    return 0;
}
//  int main() {
//     kvArr = malloc(sizeof(KeyValueArray));
//     initKeyValueArray(kvArr);
//      drv_write_int("test_int", 123);
//     drv_write_double("test_double", 3.14);
//     drv_write_string("test_string", "Hello World");
//      int intValue;
//     drv_read_int("test_int", &intValue);
//     printf("test_int: %d\n", intValue);
//      double doubleValue;
//     drv_read_double("test_double", &doubleValue);
//     printf("test_double: %lf\n", doubleValue);
//      char strValue[100];
//     drv_read_string("test_string", strValue, 100);
//     printf("test_string: %s\n", strValue);
    
//      freeKeyValueArray(kvArr);
//     free(kvArr);
//     return 0;
// }