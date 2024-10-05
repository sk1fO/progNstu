#include <iostream>
#include "linkedlist.h"
#include "queue.h"
#include "stack.h"
using namespace std;



int main(){
    list flist;
    flist.push_back("Lol");
    flist.push_back("123");
    flist.push_first("test");
    flist.push_first("test2");
    flist.remove_first();
    flist.remove_last();
    flist.remove("Lol");
    cout << flist.find("test") << endl;
    flist.print();

    queue fqueue;
    fqueue.push("12");
    fqueue.pop();
    fqueue.print();

    stack fstack;
    fstack.push("11");
    fstack.pop();
    fstack.print();


    return 0;
}