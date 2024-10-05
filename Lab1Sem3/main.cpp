#include <iostream>
using namespace std;

struct Node {
    string val;
    Node* next;
    Node(string _val) : val(_val), next(nullptr){}
};

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

int main(){
    list flist;
    flist.push_back("Lol");
    flist.push_back("123");
    flist.print();

    return 0;
}