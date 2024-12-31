#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <time.h>
 // 导入Golang库
#include "libadd.h"

 int main() {
    // 调用Golang导出函数
    clock_t start = clock();
    // 循环10000000次调用Golang导出函数
    for (int i = 0; i < 1000000; i++) {
        int result = add(1, 2);
    }
    // 获取当前时间
    clock_t end = clock();
    // 计算时间差
    double duration = (double)(end - start) / CLOCKS_PER_SEC;
    // 输出用时
    printf("Time taken: %f seconds\n", duration);
     // 输出结果
    //printf("Result: %d\n", result);
     return 0;
}