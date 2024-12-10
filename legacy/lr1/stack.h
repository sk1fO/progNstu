#include "node.h"

#include <iostream>
using namespace std;

struct stack{
    Node* first;
    Node* last;

    stack() : first(nullptr), last(nullptr){}

    bool isEmpty(){
        return first == nullptr;
    }

    void push(string _val){
        Node* p = new Node (_val);
        if(isEmpty()) {
            first = p;
            last = p;
            return;
        }
        last -> next = p;
        last = p;
    }

    void pop() {
        if (isEmpty()) return;
        if (first == last) {
            Node* p = first;
            first = p->next;
            delete p;
            return;
        }

        Node* p = first;
        while (p->next != last) p = p->next;
        p->next = nullptr;
        delete last;
        last = p;
    }


    void print(){
        if (isEmpty()) return;

        Node* p = first;
        while (p){
            cout << p -> val << " ";
            p = p->next;
        }
        cout << endl;
    }

};