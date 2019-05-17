// online c editor - https://onlinegdb.com/HySykSJoE

#include <stdio.h>

int* f() {
    int a;
    a = 10;
    return &a;
}

void main()
{
    int* p = f();
    printf("p is: %x\n", p);   // p is 0
    printf("*p is: %d\n", *p); // segmentation fault

	// 
}

