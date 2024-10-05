#pragma once

#include <iostream>

using namespace std;

struct Node {
    string val;
    Node* next;
    Node(string _val) : val(_val), next(nullptr){}
};