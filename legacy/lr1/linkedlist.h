#include "node.h"

#include <iostream>
using namespace std;

struct list{
    Node* first;
    Node* last;

    list() : first(nullptr), last(nullptr){}

    bool isEmpty(){
        return first == nullptr;
    }


    void push_back(string _val){
        Node* p = new Node (_val);
        if(isEmpty()) {
            first = p;
            last = p;
            return;
        }
        last -> next = p;
        last = p;
    }

    void push_first(string _val){
        Node* p = new Node (_val);
        if(isEmpty()) {
            first = p;
            last = p;
            return;
        }

        p->next = first;
        first = p;
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
    
    Node* find(string _val) {
        Node* p = first;
        while (p && p->val != _val) p = p->next;
        return (p && p->val == _val) ? p : nullptr;
    }

    void remove_first() {
        if (isEmpty()) return;
        Node* p = first;
        first = p->next;
        delete p;
    }

    void remove_last() {
        if (isEmpty()) return;
        if (first == last) {
            remove_first();
            return;
        }
        Node* p = first;
        while (p->next != last) p = p->next;
        p->next = nullptr;
        delete last;
        last = p;
    }

    void remove(string _val) {
        if (isEmpty()) return;
        if (first->val == _val) {
            remove_first();
            return;
        }
        else if (last->val == _val) {
            remove_last();
            return;
        }
        Node* slow = first;
        Node* fast = first->next;
        while (fast && fast->val != _val) {
            fast = fast->next;
            slow = slow->next;
        }
        if (!fast) {
            cout << "This element does not exist" << endl;
            return;
        }
        slow->next = fast->next;
        delete fast;
    }



};



